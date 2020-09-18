#!/usr/bin/env bash
cp -f ./test/testdata/k8s-e2e-test-driver.yaml /tmp/cryptpvc-csi-k8s-e2e-test-driver.yaml
cd ~/go/src/k8s.io/kubernetes && kubetest --deployment=local --test --check-version-skew=false --test_args="--ginkgo.focus=External.*Storage.* --storage.testdriver=/tmp/cryptpvc-csi-k8s-e2e-test-driver.yaml"
