package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	connect "github.com/bufbuild/connect-go"
	"golang.org/x/exp/slog"

	v1 "github.com/emahiro/grpc-go-api-skelton/gen/proto/greet/v1"
)

type GreeterService struct{}

func (s *GreeterService) Greet(ctx context.Context, req *connect.Request[v1.GreetRequest]) (*connect.Response[v1.GreetResponse], error) {
	userName := req.Msg.UserName
	if userName == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("user name is empty"))
	}
	resp := connect.NewResponse(&v1.GreetResponse{
		Message: fmt.Sprintf("Hello %s", userName),
	})
	return resp, nil
}

func (s *GreeterService) GreetStreaming(ctx context.Context, stream *connect.ClientStream[v1.GreetStreamingRequest]) (*connect.Response[v1.GreetStreamingResponse], error) {
	slog.InfoCtx(ctx, "request header", "header", stream.RequestHeader())
	var message strings.Builder
	for stream.Receive() {
		g := fmt.Sprintf("Hello, %s\n", stream.Msg().UserName)
		if _, err := message.WriteString(g); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}
	if err := stream.Err(); err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}
	res := connect.NewResponse(&v1.GreetStreamingResponse{
		Message: message.String(),
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}

func (s *GreeterService) GreetDidiStreaming(ctx context.Context, stream *connect.BidiStream[v1.GreetDidiStreamingRequest, v1.GreetDidiStreamingResponse]) error {
	return nil
}
