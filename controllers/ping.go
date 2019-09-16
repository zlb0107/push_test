package controllers

import (
	"context"
	"net/http"
)

type PingHandler struct {
}

func (*PingHandler) Serve(ctx context.Context, req *http.Request) (interface{}, int) {
	return "pong", 0
}
