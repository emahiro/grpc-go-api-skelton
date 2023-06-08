package main

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/exp/slog"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/bufbuild/connect-go"
	otelconnect "github.com/bufbuild/connect-opentelemetry-go"

	"github.com/emahiro/grpc-go-api-skelton/gen/proto/echo/v1/echov1connect"
	"github.com/emahiro/grpc-go-api-skelton/gen/proto/greet/v1/greetv1connect"
	"github.com/emahiro/grpc-go-api-skelton/intercepter"
	"github.com/emahiro/grpc-go-api-skelton/service"
)

var addr = ":8080"

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	intercepters := connect.WithInterceptors(
		intercepter.NewIntercepter(),
		otelconnect.NewInterceptor(),
	)

	mux := http.NewServeMux()
	mux.Handle(greetv1connect.NewGreetServiceHandler(&service.GreeterService{}, intercepters))
	mux.Handle(echov1connect.NewEchoServiceHandler(&service.EchoService{}, intercepters))

	server := &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	go func() {
		<-ctx.Done()
		if err := server.Close(); err != nil {
			panic(err)
		}
	}()

	slog.InfoCtx(ctx, "start server", "port", fmt.Sprintf("localhost%s", addr))
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
