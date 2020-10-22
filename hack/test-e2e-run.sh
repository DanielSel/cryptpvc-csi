#!/usr/bin/env bash
TESTDRIVER_PATH="$PWD/test/testdata/k8s-e2e-test-driver.yml"
if [ ! -d "./kubernetes" ]; then
    ./blob/kubetest --deployment=local --extract $(cat ./versions/kubernetes.txt) --test --check-version-skew=false --test_args="--ginkgo.focus=External.*Storage.* --storage.testdriver=${TESTDRIVER_PATH}"
else
    (cd ./kubernetes && ../.blob/kubetest --deployment=local ${EXTRACT_ARGS} --test --check-version-skew=false --test_args="--ginkgo.focus=External.*Storage.* --storage.testdriver=${TESTDRIVER_PATH}")
fi
