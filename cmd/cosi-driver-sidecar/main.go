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

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"
	"errors"
	"net"
	"strings"

	"google.golang.org/grpc"
	"k8s.io/klog"
 
        "github.com/brahmaroutu/cosi-external-provisioner/controller"
        "github.com/kubernetes-sigs/cosi-driver-sidecar/pkg/server"
        "github.com/container-object-storage-interface/spec/lib/go/cosi"
)

const (
	// Default timeout of short COSI calls like CreateBucket
	cosiTimeout = time.Second

	// Verify (and update, if needed) the node ID at this freqeuency.
	sleepDuration = 2 * time.Minute

	// Interval of logging connection errors
	connectionLoggingInterval = 10 * time.Second
)

// Command line flags
var (
	connectionTimeout       = flag.Duration("connection-timeout", 0, "The --connection-timeout flag is deprecated")
	cosiAddress             = flag.String("cosi-address", "/run/cosi/socket", "Path of the COSI driver socket that the provisioner  will connect to.")
	showVersion             = flag.Bool("version", false, "Show version.")
	version                 = "unknown"

	// List of supported versions
	supportedVersions = []string{"1.0.0"}

        //provisionController *controller.ProvisionController
)


type options struct {
	reconnect func() bool
}

// Option is the type of all optional parameters for Connect.
type Option func(o *options)

// LogGRPC is gPRC unary interceptor for logging of CSI messages at level 5. It removes any secrets from the message.
func LogGRPC(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	klog.V(5).Infof("GRPC call: %s", method)
	klog.V(5).Infof("GRPC request: %s", req)
	err := invoker(ctx, method, req, reply, cc, opts...)
	klog.V(5).Infof("GRPC response: %s", reply)
	klog.V(5).Infof("GRPC error: %v", err)
	return err
}

func main() {
        flag.CommandLine.Parse([]string{})
	flag.Set("logtostderr", "true")
	flag.Parse()

	if *showVersion {
		fmt.Println(os.Args[0], version)
		return
	}
	klog.Infof("Version: %s", version)

	if *connectionTimeout != 0 {
		klog.Warning("--connection-timeout is deprecated and will have no effect")
	}

        cds := server.DriverServer{"testDriver", "1.0"}
        go Serve(*cosiAddress, "default", cds)


	klog.V(1).Infof("Attempting to open a gRPC connection with: %q", *cosiAddress)
	cosiConn, err := connect(*cosiAddress, []grpc.DialOption{}, nil)
	if err != nil {
		klog.Errorf("error connecting to COSI driver: %v", err)
		os.Exit(1)
	}

	klog.V(1).Infof("Calling COSI driver to discover driver name")
	ctx, cancel := context.WithTimeout(context.Background(), cosiTimeout)
	defer cancel()

        fmt.Println("creating Client")
        client := cosi.NewCOSIDriverClient(cosiConn)
	req := cosi.CreateBucketRequest{Name: "testBucket"}
        fmt.Println("call CreateBucket")
	rsp, err := client.CreateBucket(ctx, &req)
        fmt.Println("Got response")
	if err != nil {
		klog.Errorf("error calling COSI CreateBucket: %v", err)
		fmt.Errorf("error calling COSI CreateBucket: %v", err)
		os.Exit(1)
	}

	fmt.Println("COSI CreateBucket returned : %v", rsp.Bucket)
/*
	provisionController = controller.NewProvisionController(
		clientset,
		provisionerName,
		csiProvisioner,
		serverVersion.GitVersion,
		provisionerOptions...,
	)

	run := func(context.Context) {
		stopCh := context.Background().Done()
		factory.Start(stopCh)
		cacheSyncResult := factory.WaitForCacheSync(stopCh)
		for _, v := range cacheSyncResult {
			if !v {
				klog.Fatalf("Failed to sync Informers!")
			}
		}

		provisionController.Run(wait.NeverStop)
	}
*/
}

func Serve(endpoint, ns string, cds cosi.COSIDriverServer ) {
	s := server.NewNonBlockingGRPCServer()
	s.Start(endpoint,cds)
	s.Wait()
}

// connect is the internal implementation of Connect. It has more options to enable testing.
func connect(
	address string,
	dialOptions []grpc.DialOption, connectOptions []Option) (*grpc.ClientConn, error) {
	var o options
	for _, option := range connectOptions {
		option(&o)
	}

	dialOptions = append(dialOptions,
		grpc.WithInsecure(),                   // Don't use TLS, it's usually local Unix domain socket in a container.
		grpc.WithBackoffMaxDelay(time.Second), // Retry every second after failure.
		grpc.WithBlock(),                      // Block until connection succeeds.
		grpc.WithChainUnaryInterceptor(
			LogGRPC, // Log all messages.
		),
	)
	unixPrefix := "unix://"
	if strings.HasPrefix(address, "/") {
		// It looks like filesystem path.
		address = unixPrefix + address
	}

	if strings.HasPrefix(address, unixPrefix) {
		// state variables for the custom dialer
		haveConnected := false
		lostConnection := false
		reconnect := true

		dialOptions = append(dialOptions, grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			if haveConnected && !lostConnection {
				// We have detected a loss of connection for the first time. Decide what to do...
				// Record this once. TODO (?): log at regular time intervals.
				klog.Errorf("Lost connection to %s.", address)
				// Inform caller and let it decide? Default is to reconnect.
				if o.reconnect != nil {
					reconnect = o.reconnect()
				}
				lostConnection = true
			}
			if !reconnect {
				return nil, errors.New("connection lost, reconnecting disabled")
			}
			conn, err := net.DialTimeout("unix", address[len(unixPrefix):], timeout)
			if err == nil {
				// Connection restablished.
				haveConnected = true
				lostConnection = false
			}
			return conn, err
		}))
	} else if o.reconnect != nil {
		return nil, errors.New("OnConnectionLoss callback only supported for unix:// addresses")
	}

	klog.Infof("Connecting to %s", address)

	// Connect in background.
	var conn *grpc.ClientConn
	var err error
	ready := make(chan bool)
	go func() {
		conn, err = grpc.Dial(address, dialOptions...)
		close(ready)
	}()

	// Log error every connectionLoggingInterval
	ticker := time.NewTicker(connectionLoggingInterval)
	defer ticker.Stop()

	// Wait until Dial() succeeds.
	for {
		select {
		case <-ticker.C:
			klog.Warningf("Still connecting to %s", address)

		case <-ready:
			return conn, err
		}
	}
}

