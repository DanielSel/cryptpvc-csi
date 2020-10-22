#!/usr/bin/env bash
cp -f ./test/testdata/k8s-e2e-test-driver.yaml /tmp/cryptpvc-csi-k8s-e2e-test-driver.yaml
.blob/kubetest --deployment=local --extract $(cat ./versions/kubernetes.txt) --test --check-version-skew=false --test_args="--ginkgo.focus=External.*Storage.* --storage.testdriver=/tmp/cryptpvc-csi-k8s-e2e-test-driver.yaml"
