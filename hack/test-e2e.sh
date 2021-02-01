#!/usr/bin/env bash

kind export kubeconfig --name cryptpvc-e2e
./hack/deploy-snapshot-crds.sh
./hack/deploy.sh
./hack/test-e2e-run.sh
