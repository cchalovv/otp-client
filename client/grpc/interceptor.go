package grpc_client

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/cchalovv/otp-client/pkg/errs"
	"github.com/cchalovv/otp-client/pkg/proto/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// error interceptor
type grpcClientInterceptorErrorT struct {
	errMessagePrefix string
}

func (o *grpcClientInterceptorErrorT) grpcClientInterceptorError(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if ok {
		if len(st.Details()) > 0 {
			stDetail := st.Details()[0]
			errObj, ok := stDetail.(*common.ErrorRep)
			if ok {
				return errs.ErrFull{
					Err:    errs.Err(errObj.Code),
					Desc:   o.errMessagePrefix + errObj.Message,
					Fields: errObj.Fields,
				}
			}
		}
		return errs.ErrFull{
			Err:  errs.ServiceNA,
			Desc: o.errMessagePrefix + st.String(),
		}
	}

	return fmt.Errorf(o.errMessagePrefix+": %w", err)
}

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
