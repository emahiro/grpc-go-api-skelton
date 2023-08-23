package main

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	protovalidate "github.com/bufbuild/protovalidate-go"
	"golang.org/x/exp/slog"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	ddotel "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentelemetry"

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

	ddTraceProvider := ddotel.NewTracerProvider()
	defer func() {
		if err := ddTraceProvider.Shutdown(); err != nil {
			panic(err)
		}
	}()
	intercepters := connect.WithInterceptors(
		intercepter.NewIntercepter(),
		otelconnect.NewInterceptor(
			otelconnect.WithTracerProvider(ddTraceProvider), // Set custom tracer provider
		),
		intercepter.NewAuthIntercepter(),
	)

	v, err := protovalidate.New()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle(greetv1connect.NewGreetServiceHandler(&service.GreeterService{
		Validator: v,
	}, intercepters))
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

	slog.InfoCtx(ctx, "start server", "port", "localhost"+addr)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
