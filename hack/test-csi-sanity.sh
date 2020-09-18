#!/bin/bash

## Must be run from the root of the repo

UDS="/tmp/e2e-csi-sanity.sock"
CSI_ENDPOINT="unix://${UDS}"
CSI_MOUNTPOINT="/tmp/e2e-csi"
APP=cryptpvc-csi

SKIP="WithCapacity"


# Cleanup
rm -f $UDS

# Ensure csi-sanity
which csi-sanity
if [ $? -ne 0 ]; then
    echo "Installing csi-sanity..."
    curl -Lo csi-sanity.go https://raw.githubusercontent.com/kubernetes-csi/csi-test/master/cmd/csi-sanity/main.go &&\
    go build -o $GOPATH/bin/csi-sanity csi-sanity.go
fi

# Build cryptpvc-csi
go build -o ./bin/cryptpvc-csi ./cmd/cryptpvc-csi

# Start the application in the background
./bin/$APP --endpoint=$CSI_ENDPOINT --nodeid=1 &
pid=$!

# Need to skip Capacity testing since cryptpvc-csi does not support it
csi-sanity $@ \
    --ginkgo.skip=${SKIP} \
    --csi.mountdir=$CSI_MOUNTPOINT \
    --csi.endpoint=$CSI_ENDPOINT ; ret=$?
kill -9 $pid
rm -f $UDS

if [ $ret -ne 0 ] ; then
	exit $ret
fi

exit 0
