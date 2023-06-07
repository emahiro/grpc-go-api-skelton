package service

import (
	"context"

	connect_go "github.com/bufbuild/connect-go"

	v1 "github.com/emahiro/grpc-go-api-skelton/gen/proto/echo/v1"
)

type EchoService struct{}

func (s *EchoService) Echo(ctx context.Context, req *connect_go.Request[v1.EchoRequest]) (*connect_go.Response[v1.EchoResponse], error) {
	resp := connect_go.NewResponse(&v1.EchoResponse{
		Message: req.Msg.Message,
	})
	return resp, nil
}
