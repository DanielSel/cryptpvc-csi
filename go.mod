module github.com/ds-os/cryptpvc-csi

go 1.15

require (
	github.com/Rican7/retry v0.1.0 // indirect
	github.com/bazelbuild/bazel-gazelle v0.19.1-0.20191105222053-70208cbdc798 // indirect
	github.com/cespare/prettybench v0.0.0-20150116022406-03b8cfe5406c // indirect
	github.com/checkpoint-restore/go-criu v0.0.0-20181120144056-17b0214f6c48 // indirect
	github.com/container-storage-interface/spec v1.3.0
	github.com/godbus/dbus v0.0.0-20181101234600-2ff6f7ffd60f // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.2 // indirect
	github.com/kubernetes-csi/csi-lib-utils v0.3.0
	github.com/mattn/go-shellwords v1.0.5 // indirect
	github.com/mesos/mesos-go v0.0.9 // indirect
	github.com/pborman/uuid v1.2.0
	github.com/pquerna/ffjson v0.0.0-20180717144149-af8b230fcd20 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/stretchr/testify v1.4.0 // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	google.golang.org/grpc v1.29.1
	gotest.tools/gotestsum v0.3.5 // indirect
	k8s.io/kubernetes v1.19.1
	k8s.io/repo-infra v0.0.1-alpha.1 // indirect
	k8s.io/utils v0.0.0-20200729134348-d5654de09c73
	sigs.k8s.io/structured-merge-diff/v3 v3.0.0-20200116222232-67a7b8c61874 // indirect
)

replace k8s.io/api => k8s.io/api v0.19.1
replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.1
replace k8s.io/apimachinery => k8s.io/apimachinery v0.19.4-rc.0
replace k8s.io/apiserver => k8s.io/apiserver v0.19.1
replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.19.1
replace k8s.io/client-go => k8s.io/client-go v0.19.1
replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.19.1
replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.19.1
replace k8s.io/code-generator => k8s.io/code-generator v0.19.4-rc.0
replace k8s.io/component-base => k8s.io/component-base v0.19.1
replace k8s.io/cri-api => k8s.io/cri-api v0.19.4-rc.0
replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.19.1
replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.19.1
replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.19.1
replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.19.1
replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.19.1
replace k8s.io/kubectl => k8s.io/kubectl v0.19.1
replace k8s.io/kubelet => k8s.io/kubelet v0.19.1
replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.19.1
replace k8s.io/metrics => k8s.io/metrics v0.19.1
replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.19.1
replace k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.19.1
replace k8s.io/sample-controller => k8s.io/sample-controller v0.19.1
replace k8s.io/controller-manager => k8s.io/controller-manager v0.19.4-rc.0
