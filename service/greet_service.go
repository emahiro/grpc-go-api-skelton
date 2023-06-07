package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"
	connect_go "github.com/bufbuild/connect-go"

	v1 "github.com/emahiro/grpc-go-api-skelton/gen/proto/greet/v1"
)

type GreeterService struct{}

func (s *GreeterService) Greet(ctx context.Context, req *connect_go.Request[v1.GreetRequest]) (*connect_go.Response[v1.GreetResponse], error) {
	userName := req.Msg.UserName
	if userName == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("user name is empty"))
	}
	resp := connect_go.NewResponse(&v1.GreetResponse{
		Message: fmt.Sprintf("Hello %s", userName),
	})
	return resp, nil
}
