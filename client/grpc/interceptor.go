package grpc_client

import (
	"context"
	"encoding/base64"
)

type grpcClientBasicAuth struct {
	metadata map[string]string
}

func newGrpcClientBasicAuth(username, password string) *grpcClientBasicAuth {
	return &grpcClientBasicAuth{
		metadata: map[string]string{
			"authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password)),
		},
	}
}

func (o *grpcClientBasicAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return o.metadata, nil
}

func (o *grpcClientBasicAuth) RequireTransportSecurity() bool {
	return false
}
