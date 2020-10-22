#!/usr/bin/env bash
.blob/kubetest --deployment=local --extract $(cat ./versions/kubernetes.txt) --test --check-version-skew=false --test_args="--ginkgo.focus=External.*Storage.* --storage.testdriver=./test/testdata/k8s-e2e-test-driver.yaml"
