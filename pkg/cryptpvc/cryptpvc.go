/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cryptpvc

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/kubernetes/pkg/volume/util/volumepathhandler"
	utilexec "k8s.io/utils/exec"

	timestamp "github.com/golang/protobuf/ptypes/timestamp"
)

const (
	kib    int64 = 1024
	mib    int64 = kib * 1024
	gib    int64 = mib * 1024
	gib100 int64 = gib * 100
	tib    int64 = gib * 1024
	tib100 int64 = tib * 100
)

type cryptpvc struct {
	name              string
	nodeID            string
	version           string
	endpoint          string
	ephemeral         bool
	maxVolumesPerNode int64

	ids *identityServer
	ns  *nodeServer
	cs  *controllerServer
}

type cryptpvcVolume struct {
	VolName       string     `json:"volName"`
	VolID         string     `json:"volID"`
	VolSize       int64      `json:"volSize"`
	VolPath       string     `json:"volPath"`
	VolAccessType accessType `json:"volAccessType"`
	ParentVolID   string     `json:"parentVolID,omitempty"`
	ParentSnapID  string     `json:"parentSnapID,omitempty"`
	Ephemeral     bool       `json:"ephemeral"`
}

type cryptpvcSnapshot struct {
	Name         string              `json:"name"`
	Id           string              `json:"id"`
	VolID        string              `json:"volID"`
	Path         string              `json:"path"`
	CreationTime timestamp.Timestamp `json:"creationTime"`
	SizeBytes    int64               `json:"sizeBytes"`
	ReadyToUse   bool                `json:"readyToUse"`
}

var (
	vendorVersion = "dev"

	cryptpvcVolumes         map[string]cryptpvcVolume
	cryptpvcVolumeSnapshots map[string]cryptpvcSnapshot
)

const (
	// Directory where data for volumes and snapshots are persisted.
	// This can be ephemeral within the container or persisted if
	// backed by a Pod volume.
	dataRoot = "/csi-data-dir"

	// Extension with which snapshot files will be saved.
	snapshotExt = ".snap"
)

func init() {
	cryptpvcVolumes = map[string]cryptpvcVolume{}
	cryptpvcVolumeSnapshots = map[string]cryptpvcSnapshot{}
}

func NewCryptpvcDriver(driverName, nodeID, endpoint string, ephemeral bool, maxVolumesPerNode int64, version string) (*cryptpvc, error) {
	if driverName == "" {
		return nil, errors.New("no driver name provided")
	}

	if nodeID == "" {
		return nil, errors.New("no node id provided")
	}

	if endpoint == "" {
		return nil, errors.New("no driver endpoint provided")
	}
	if version != "" {
		vendorVersion = version
	}

	if err := os.MkdirAll(dataRoot, 0750); err != nil {
		return nil, fmt.Errorf("failed to create dataRoot: %v", err)
	}

	glog.Infof("Driver: %v ", driverName)
	glog.Infof("Version: %s", vendorVersion)

	return &cryptpvc{
		name:              driverName,
		version:           vendorVersion,
		nodeID:            nodeID,
		endpoint:          endpoint,
		ephemeral:         ephemeral,
		maxVolumesPerNode: maxVolumesPerNode,
	}, nil
}

func getSnapshotID(file string) (bool, string) {
	glog.V(4).Infof("file: %s", file)
	// Files with .snap extension are volumesnapshot files.
	// e.g. foo.snap, foo.bar.snap
	if filepath.Ext(file) == snapshotExt {
		return true, strings.TrimSuffix(file, snapshotExt)
	}
	return false, ""
}

func discoverExistingSnapshots() {
	glog.V(4).Infof("discovering existing snapshots in %s", dataRoot)
	files, err := ioutil.ReadDir(dataRoot)
	if err != nil {
		glog.Errorf("failed to discover snapshots under %s: %v", dataRoot, err)
	}
	for _, file := range files {
		isSnapshot, snapshotID := getSnapshotID(file.Name())
		if isSnapshot {
			glog.V(4).Infof("adding snapshot %s from file %s", snapshotID, getSnapshotPath(snapshotID))
			cryptpvcVolumeSnapshots[snapshotID] = cryptpvcSnapshot{
				Id:         snapshotID,
				Path:       getSnapshotPath(snapshotID),
				ReadyToUse: true,
			}
		}
	}
}

func (hp *cryptpvc) Run() {
	// Create GRPC servers
	hp.ids = NewIdentityServer(hp.name, hp.version)
	hp.ns = NewNodeServer(hp.nodeID, hp.ephemeral, hp.maxVolumesPerNode)
	hp.cs = NewControllerServer(hp.ephemeral, hp.nodeID)

	discoverExistingSnapshots()
	s := NewNonBlockingGRPCServer()
	s.Start(hp.endpoint, hp.ids, hp.cs, hp.ns)
	s.Wait()
}

func getVolumeByID(volumeID string) (cryptpvcVolume, error) {
	if cryptpvcVol, ok := cryptpvcVolumes[volumeID]; ok {
		return cryptpvcVol, nil
	}
	return cryptpvcVolume{}, fmt.Errorf("volume id %s does not exist in the volumes list", volumeID)
}

func getVolumeByName(volName string) (cryptpvcVolume, error) {
	for _, cryptpvcVol := range cryptpvcVolumes {
		if cryptpvcVol.VolName == volName {
			return cryptpvcVol, nil
		}
	}
	return cryptpvcVolume{}, fmt.Errorf("volume name %s does not exist in the volumes list", volName)
}

func getSnapshotByName(name string) (cryptpvcSnapshot, error) {
	for _, snapshot := range cryptpvcVolumeSnapshots {
		if snapshot.Name == name {
			return snapshot, nil
		}
	}
	return cryptpvcSnapshot{}, fmt.Errorf("snapshot name %s does not exist in the snapshots list", name)
}

// getVolumePath returns the canonical path for cryptpvc volume
func getVolumePath(volID string) string {
	return filepath.Join(dataRoot, volID)
}

// createVolume create the directory for the cryptpvc volume.
// It returns the volume path or err if one occurs.
func createCryptpvcVolume(volID, name string, cap int64, volAccessType accessType, ephemeral bool) (*cryptpvcVolume, error) {
	path := getVolumePath(volID)

	switch volAccessType {
	case mountAccess:
		err := os.MkdirAll(path, 0777)
		if err != nil {
			return nil, err
		}
	case blockAccess:
		executor := utilexec.New()
		size := fmt.Sprintf("%dM", cap/mib)
		// Create a block file.
		out, err := executor.Command("fallocate", "-l", size, path).CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("failed to create block device: %v, %v", err, string(out))
		}

		// Associate block file with the loop device.
		volPathHandler := volumepathhandler.VolumePathHandler{}
		_, err = volPathHandler.AttachFileDevice(path)
		if err != nil {
			// Remove the block file because it'll no longer be used again.
			if err2 := os.Remove(path); err2 != nil {
				glog.Errorf("failed to cleanup block file %s: %v", path, err2)
			}
			return nil, fmt.Errorf("failed to attach device %v: %v", path, err)
		}
	default:
		return nil, fmt.Errorf("unsupported access type %v", volAccessType)
	}

	cryptpvcVol := cryptpvcVolume{
		VolID:         volID,
		VolName:       name,
		VolSize:       cap,
		VolPath:       path,
		VolAccessType: volAccessType,
		Ephemeral:     ephemeral,
	}
	cryptpvcVolumes[volID] = cryptpvcVol
	return &cryptpvcVol, nil
}

// updateVolume updates the existing cryptpvc volume.
func updateCryptpvcVolume(volID string, volume cryptpvcVolume) error {
	glog.V(4).Infof("updating cryptpvc volume: %s", volID)

	if _, err := getVolumeByID(volID); err != nil {
		return err
	}

	cryptpvcVolumes[volID] = volume
	return nil
}

// deleteVolume deletes the directory for the cryptpvc volume.
func deleteCryptpvcVolume(volID string) error {
	glog.V(4).Infof("deleting cryptpvc volume: %s", volID)

	vol, err := getVolumeByID(volID)
	if err != nil {
		// Return OK if the volume is not found.
		return nil
	}

	if vol.VolAccessType == blockAccess {
		volPathHandler := volumepathhandler.VolumePathHandler{}
		// Get the associated loop device.
		device, err := volPathHandler.GetLoopDevice(getVolumePath(volID))
		if err != nil {
			return fmt.Errorf("failed to get the loop device: %v", err)
		}

		if device != "" {
			// Remove any associated loop device.
			glog.V(4).Infof("deleting loop device %s", device)
			if err := volPathHandler.DetachFileDevice(device); err != nil {
				return fmt.Errorf("failed to remove loop device %v: %v", device, err)
			}
		}
	}

	path := getVolumePath(volID)
	if err := os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	delete(cryptpvcVolumes, volID)
	return nil
}

// cryptpvcIsEmpty is a simple check to determine if the specified cryptpvc directory
// is empty or not.
func cryptpvcIsEmpty(p string) (bool, error) {
	f, err := os.Open(p)
	if err != nil {
		return true, fmt.Errorf("unable to open cryptpvc volume, error: %v", err)
	}
	defer f.Close()

	_, err = f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

// loadFromSnapshot populates the given destPath with data from the snapshotID
func loadFromSnapshot(size int64, snapshotId, destPath string, mode accessType) error {
	snapshot, ok := cryptpvcVolumeSnapshots[snapshotId]
	if !ok {
		return status.Errorf(codes.NotFound, "cannot find snapshot %v", snapshotId)
	}
	if snapshot.ReadyToUse != true {
		return status.Errorf(codes.Internal, "snapshot %v is not yet ready to use.", snapshotId)
	}
	if snapshot.SizeBytes > size {
		return status.Errorf(codes.InvalidArgument, "snapshot %v size %v is greater than requested volume size %v", snapshotId, snapshot.SizeBytes, size)
	}
	snapshotPath := snapshot.Path

	var cmd []string
	switch mode {
	case mountAccess:
		cmd = []string{"tar", "zxvf", snapshotPath, "-C", destPath}
	case blockAccess:
		cmd = []string{"dd", "if=" + snapshotPath, "of=" + destPath}
	default:
		return status.Errorf(codes.InvalidArgument, "unknown accessType: %d", mode)
	}
	executor := utilexec.New()
	out, err := executor.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if err != nil {
		return status.Errorf(codes.Internal, "failed pre-populate data from snapshot %v: %v: %s", snapshotId, err, out)
	}
	return nil
}

// loadFromVolume populates the given destPath with data from the srcVolumeID
func loadFromVolume(size int64, srcVolumeId, destPath string, mode accessType) error {
	cryptpvcVolume, ok := cryptpvcVolumes[srcVolumeId]
	if !ok {
		return status.Error(codes.NotFound, "source volumeId does not exist, are source/destination in the same storage class?")
	}
	if cryptpvcVolume.VolSize > size {
		return status.Errorf(codes.InvalidArgument, "volume %v size %v is greater than requested volume size %v", srcVolumeId, cryptpvcVolume.VolSize, size)
	}
	if mode != cryptpvcVolume.VolAccessType {
		return status.Errorf(codes.InvalidArgument, "volume %v mode is not compatible with requested mode", srcVolumeId)
	}

	switch mode {
	case mountAccess:
		return loadFromFilesystemVolume(cryptpvcVolume, destPath)
	case blockAccess:
		return loadFromBlockVolume(cryptpvcVolume, destPath)
	default:
		return status.Errorf(codes.InvalidArgument, "unknown accessType: %d", mode)
	}
}

func loadFromFilesystemVolume(cryptpvcVolume cryptpvcVolume, destPath string) error {
	srcPath := cryptpvcVolume.VolPath
	isEmpty, err := cryptpvcIsEmpty(srcPath)
	if err != nil {
		return status.Errorf(codes.Internal, "failed verification check of source cryptpvc volume %v: %v", cryptpvcVolume.VolID, err)
	}

	// If the source cryptpvc volume is empty it's a noop and we just move along, otherwise the cp call will fail with a a file stat error DNE
	if !isEmpty {
		args := []string{"-a", srcPath + "/.", destPath + "/"}
		executor := utilexec.New()
		out, err := executor.Command("cp", args...).CombinedOutput()
		if err != nil {
			return status.Errorf(codes.Internal, "failed pre-populate data from volume %v: %v: %s", cryptpvcVolume.VolID, err, out)
		}
	}
	return nil
}

func loadFromBlockVolume(cryptpvcVolume cryptpvcVolume, destPath string) error {
	srcPath := cryptpvcVolume.VolPath
	args := []string{"if=" + srcPath, "of=" + destPath}
	executor := utilexec.New()
	out, err := executor.Command("dd", args...).CombinedOutput()
	if err != nil {
		return status.Errorf(codes.Internal, "failed pre-populate data from volume %v: %v: %s", cryptpvcVolume.VolID, err, out)
	}
	return nil
}
