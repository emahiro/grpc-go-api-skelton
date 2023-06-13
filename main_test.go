package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/emahiro/grpc-go-api-skelton/gen/proto/greet/v1/greetv1connect"
	"github.com/emahiro/grpc-go-api-skelton/service"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestGreetService_Greet(t *testing.T) {
	path, handler := greetv1connect.NewGreetServiceHandler(&service.GreeterService{})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodPost, ts.URL+path+"Greet", bytes.NewBuffer([]byte(`{"user_name": "taro"}`)))
	req.Header.Set("Content-Type", "application/json")
	cli := &http.Client{
		Transport: ts.Client().Transport,
	}

	resp, err := cli.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status code %d, but got %d. err: %v", http.StatusOK, resp.StatusCode, string(b))
	}
}
