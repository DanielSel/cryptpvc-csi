#!/usr/bin/env bash

kind export kubeconfig --name cryptpvc-e2e
./hack/deploy.sh
./hack/test-e2e-run.sh
