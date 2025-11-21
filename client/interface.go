package otp

import (
	"context"
	"github.com/cchalovv/otp-client/model"
)

type Client interface {
	Connect() error
	Close() error

	Generate(ctx context.Context, req model.CreateRequest) (code string, err error)
	Verify(ctx context.Context, req model.VerifyRequest) error
}
