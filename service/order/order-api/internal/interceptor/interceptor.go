package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type CtxKey string

const (
	CtxKeyAdminID CtxKey = "adminID"
)

func HInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	adminID := ctx.Value(CtxKeyAdminID).(string)
	md := metadata.Pairs(
		"key1", "val1",
		"Key2", "val2",
		"requestID", "12345",
		"token", "mall-order-h",
		"adminID", adminID,
	)
	ctx = metadata.NewOutgoingContext(ctx, md)
	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}
