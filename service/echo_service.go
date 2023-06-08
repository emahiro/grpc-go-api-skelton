package service

import (
	"context"

	connect "github.com/bufbuild/connect-go"

	v1 "github.com/emahiro/grpc-go-api-skelton/gen/proto/echo/v1"
)

type EchoService struct{}

func (s *EchoService) Echo(ctx context.Context, req *connect.Request[v1.EchoRequest]) (*connect.Response[v1.EchoResponse], error) {
	resp := connect.NewResponse(&v1.EchoResponse{
		Message: req.Msg.Message,
	})
	return resp, nil
}
