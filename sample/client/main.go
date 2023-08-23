package main

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"golang.org/x/exp/slog"

	greetv1 "github.com/emahiro/grpc-go-api-skelton/gen/proto/greet/v1"
	"github.com/emahiro/grpc-go-api-skelton/gen/proto/greet/v1/greetv1connect"
)

func main() {
	client := greetv1connect.NewGreetServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)
	req := connect.NewRequest[greetv1.GreetRequest](&greetv1.GreetRequest{UserName: "aa"})
	req.Header().Add("Acme-Token", "test")
	res, err := client.Greet(
		context.Background(),
		req,
	)
	if err != nil {
		slog.Error("failed to call greet", "err", err)
		return
	}
	slog.Info("success to call greet", "res", res.Msg.Message)
}
