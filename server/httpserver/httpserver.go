package httpserver

import (
	"context"
	"runtime"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

func httpServer() error {
	ctxr := context.Background()
	ctx, cancel := context.WithCancel(ctxr)
	defer cancel()
	mux := runtime.NewServeMux()
}
