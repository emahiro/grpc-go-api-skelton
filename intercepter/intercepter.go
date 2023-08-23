package intercepter

import (
	"context"

	"connectrpc.com/connect"
	"golang.org/x/exp/slog"
)

func NewIntercepter() connect.UnaryInterceptorFunc {
	intercepter := func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			slog.InfoCtx(ctx, "this is intercepter")
			return next(ctx, req)
		}
	}
	return connect.UnaryInterceptorFunc(intercepter)
}
