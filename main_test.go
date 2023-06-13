package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emahiro/grpc-go-api-skelton/gen/proto/greet/v1/greetv1connect"
	"github.com/emahiro/grpc-go-api-skelton/service"
)

var path, baseURL string

func TestMain(m *testing.M) {
	p, handler := greetv1connect.NewGreetServiceHandler(&service.GreeterService{})
	path = p
	ts := httptest.NewServer(handler)
	baseURL = ts.URL
	reset := http.DefaultClient.Transport
	http.DefaultClient.Transport = ts.Client().Transport
	m.Run()
	http.DefaultTransport = reset
	ts.Close()
}

func TestGreetService_Greet(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, baseURL+path+"Greet", bytes.NewBuffer([]byte(`{"user_name": "taro"}`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status code %d, but got %d. err: %v", http.StatusOK, resp.StatusCode, string(b))
	}
}
