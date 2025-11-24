package grpc_client

import (
	"context"
	"fmt"
	otp "github.com/cchalovv/otp-client/client"
	"github.com/cchalovv/otp-client/model"
	otpPb "github.com/cchalovv/otp-client/pkg/proto/otp"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"math"
)

const clientName = "otp-client"

type Client struct {
	uri          string
	secure       bool
	username     string
	password     string
	conn         *grpc.ClientConn
	client       otpPb.OtpClient
	interceptors []grpc.DialOption
}

func NewClient(uri string, secure bool, username, password string, interceptors ...grpc.DialOption) otp.Client {
	return &Client{
		uri:          uri,
		secure:       secure,
		username:     username,
		password:     password,
		interceptors: interceptors,
	}
}

func (c *Client) Connect() (err error) {
	if c.uri == "" {
		return fmt.Errorf("otp-client uri is empty")
	}

	dialOptions := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(math.MaxInt32)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	}

	dialOptions = append(dialOptions, c.interceptors...)

	if c.secure {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	} else {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if c.username != "" {
		dialOptions = append(dialOptions, grpc.WithPerRPCCredentials(
			newGrpcClientBasicAuth(c.username, c.password),
		))
	}

	c.conn, err = grpc.NewClient(c.uri, dialOptions...)
	if err != nil {
		return fmt.Errorf("grpc.NewClient: %w", err)
	}

	c.client = otpPb.NewOtpClient(c.conn)

	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Generate(ctx context.Context, req model.GenerateRequest) (string, error) {
	resp, err := c.client.Generate(ctx, &otpPb.GenerateReq{
		Data: req.Data,
	})
	if err != nil {
		return "", err
	}

	return resp.Code, nil
}

func (c *Client) Verify(ctx context.Context, req model.VerifyRequest) error {
	_, err := c.client.Verify(ctx, &otpPb.VerifyReq{
		Data: req.Data,
		Code: req.Code,
	})
	return err
}
