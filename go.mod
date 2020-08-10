module github.com/brahmaroutu/cosi-driver-sidecar

go 1.14

require (
	github.com/brahmaroutu/cosi-external-provisioner v1.0.1
	github.com/container-object-storage-interface/api v0.0.0-20200708183033-b21b31b712bd
	github.com/container-object-storage-interface/spec v0.0.0-20200622154246-bc84d8cb63a1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/kubernetes-csi/csi-lib-utils v0.7.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae
	google.golang.org/grpc v1.30.0
	k8s.io/apimachinery v0.18.5
	k8s.io/client-go v0.18.4
	k8s.io/klog v1.0.0
)

replace github.com/brahmaroutu/cosi-external-provisioner => /home/srinib/go/src/github.com/brahmaroutu/cosi-external-provisioner

replace github.com/container-object-storage-interface/api => /home/srinib/go/src/github.com/container-object-storage-interface/api
