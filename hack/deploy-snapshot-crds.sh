# This little helper script deploys the VolumeSnapshot CRD's for kubernetes distributions that don't include them by default (e.g. kind)
# See: https://kubernetes-csi.github.io/docs/snapshot-restore-feature.html

kubectl create -f  https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/release-2.1/config/crd/snapshot.storage.k8s.io_volumesnapshotclasses.yaml
kubectl create -f  https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/release-2.1/config/crd/snapshot.storage.k8s.io_volumesnapshotcontents.yaml
kubectl create -f  https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/release-2.1/config/crd/snapshot.storage.k8s.io_volumesnapshots.yaml

# Install Snapshot Controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/release-2.1/deploy/kubernetes/snapshot-controller/rbac-snapshot-controller.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/release-2.1/deploy/kubernetes/snapshot-controller/setup-snapshot-controller.yaml