package intercepter

import (
	"context"
	"errors"

	"connectrpc.com/connect"
)

const tokenHeader = "Acme-Token"

var errNoToken = errors.New("no token provided")

type authIntercepter struct{}

func NewAuthIntercepter() *authIntercepter {
	return &authIntercepter{}
}

func (i *authIntercepter) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if req.Spec().IsClient {
			req.Header().Set(tokenHeader, "test")
		} else if req.Header().Get(tokenHeader) == "" {
			return nil, connect.NewError(connect.CodeUnauthenticated, errNoToken)
		}
		return next(ctx, req)
	})
}

func (*authIntercepter) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return connect.StreamingClientFunc(func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		conn := next(ctx, spec)
		conn.RequestHeader().Set(tokenHeader, "test")
		return conn
	})
}

func (*authIntercepter) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return connect.StreamingHandlerFunc(func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		if conn.RequestHeader().Get(tokenHeader) == "" {
			return connect.NewError(connect.CodeUnauthenticated, errNoToken)
		}
		return next(ctx, conn)
	})
}
