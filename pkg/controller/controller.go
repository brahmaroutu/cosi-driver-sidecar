/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
    "context"
	"time"
        "fmt"


	"github.com/brahmaroutu/cosi-external-provisioner/controller"
        "github.com/container-object-storage-interface/api/apis/cosi.sigs.k8s.io/v1alpha1"

 	_ "k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
//	"k8s.io/klog"

	"google.golang.org/grpc"
//	corelisters "k8s.io/client-go/listers/core/v1"
	
	"github.com/container-object-storage-interface/spec/lib/go/cosi"
        cosiclient  "github.com/container-object-storage-interface/api/clientset"


//	storagelistersv1 "k8s.io/client-go/listers/storage/v1"
//	storagelistersv1beta1 "k8s.io/client-go/listers/storage/v1beta1"
)

//secretParamsMap provides a mapping of current as well as deprecated secret keys
type secretParamsMap struct {
	name                         string
	secretNameKey                string
	secretNamespaceKey           string
}

const (
	// COSI Parameters prefixed with cosiParameterPrefix are not passed through
    // to the driver on CreateBucket calls. Instead they are intended
    // to used by the COSI external-provisioner and maybe used to populate
    // fields in subsequent COSI calls or Kubernetes API objects.
    cosiParameterPrefix = "cosi.storage.k8s.io/"

	provisionerSecretNameKey      = "csiProvisionerSecretName"
	provisionerSecretNamespaceKey = "csiProvisionerSecretNamespace"

	controllerPublishSecretNameKey      = "csiControllerPublishSecretName"
	controllerPublishSecretNamespaceKey = "csiControllerPublishSecretNamespace"

    prefixedProvisionerSecretNameKey      = cosiParameterPrefix + "provisioner-secret-name"
    prefixedProvisionerSecretNamespaceKey = cosiParameterPrefix + "provisioner-secret-namespace"

	// Bucket and BucketRequest metadata, used for sending to drivers in the  create requests, added as parameters, optional.
	bucketRequestNameKey      = "csi.storage.k8s.io/pvc/name"
	bucketRequestNamespaceKey = "csi.storage.k8s.io/pvc/namespace"
	bucketNameKey             = "csi.storage.k8s.io/pv/name"

	// Defines parameters for ExponentialBackoff used for executing
	// CSI CreateVolume API call, it gives approx 4 minutes for the CSI
	// driver to complete a volume creation.
	backoffDuration = time.Second * 5
	backoffFactor   = 1.2
	backoffSteps    = 10


	tokenBucketNameKey             = "pv.name"
	tokenBucketRequestNameKey      = "pvc.name"
	tokenBucketRequestNameSpaceKey = "pvc.namespace"

	ResyncPeriodOfCOSINodeInformer = 1 * time.Hour

	deleteVolumeRetryCount = 5

	annStorageProvisioner = "bucket.beta.kubernetes.io/cosi-provisioner"

)

var (
	provisionerSecretParams = secretParamsMap{
		name:                         "Provisioner",
		secretNameKey:                prefixedProvisionerSecretNameKey,
		secretNamespaceKey:           prefixedProvisionerSecretNamespaceKey,
	}

)

// COSIProvisioner struct
type cosiProvisioner struct {
	client                                kubernetes.Interface
        cosiInterface                         cosiclient.Interface
	cosiClient                            cosi.COSIDriverClient
	grpcClient                            *grpc.ClientConn
	timeout                               time.Duration
	identity                              string
	config                                *rest.Config
	driverName                            string
	//pluginCapabilities                    rpc.PluginCapabilitySet
	//controllerCapabilities                rpc.ControllerCapabilitySet
}

var _ controller.Provisioner = &cosiProvisioner{}

var (
	// Each provisioner have a identify string to distinguish with others. This
	// identify string will be added in Bucket annoations under this key.
	provisionerIDKey = "storage.kubernetes.io/cosiProvisionerIdentity"
)

func GetDriverName(conn *grpc.ClientConn, timeout time.Duration) (string, error) {
    //srini: this is required once we define this api.
    return "testDriver", nil
}

//func GetDriverCapabilities(conn *grpc.ClientConn, timeout time.Duration) (rpc.PluginCapabilitySet, rpc.ControllerCapabilitySet, error) {
   //srini : required to get capabilities of controller and plugin
//}

// NewCOSIProvisioner creates new CSI provisioner
func NewCOSIProvisioner(client kubernetes.Interface,
        client_cosi cosiclient.Interface,
	connectionTimeout time.Duration,
	identity string,
	grpcClient *grpc.ClientConn,
	driverName string,
//	pluginCapabilities rpc.PluginCapabilitySet,
//	controllerCapabilities rpc.ControllerCapabilitySet,
) controller.Provisioner {

    fmt.Println("create cosi client")
    cosiClient := cosi.NewCOSIDriverClient(grpcClient)

    fmt.Println("create provisioner")
	provisioner := &cosiProvisioner{
		client:                                client,
                cosiInterface:                         client_cosi,
		grpcClient:                            grpcClient,
		cosiClient:                            cosiClient,
		timeout:                               connectionTimeout,
		identity:                              identity,
		driverName:                            driverName,
//		pluginCapabilities:                    pluginCapabilities,
//		controllerCapabilities:                controllerCapabilities,
	}
	return provisioner
}

func (p *cosiProvisioner) Provision(context.Context, controller.ProvisionOptions) (*v1alpha1.Bucket, controller.ProvisioningState, error) {
    return nil, "", nil	
}

