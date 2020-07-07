package server

import (
	"github.com/container-object-storage-interface/spec/lib/go/cosi"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DriverServer struct {
     Name, Version string
}

func (ds DriverServer) CreateBucket(ctx context.Context, req *cosi.CreateBucketRequest) (*cosi.CreateBucketResponse, error) {
	glog.V(5).Infof("Using mocked CreateBucket call")

	if ds.Name == "" {
		return nil, status.Error(codes.Unavailable, "Driver name not configured")
	}

	if ds.Version == "" {
		return nil, status.Error(codes.Unavailable, "Driver is missing version")
	}

        
	return &cosi.CreateBucketResponse{
		Bucket:       &cosi.Bucket{BucketId: "1111111111" },
	}, nil
}

