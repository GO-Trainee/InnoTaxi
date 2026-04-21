package http

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbentity "awesomeProject/shared/proto/service_name"
)

// HttpHandler combines the gRPC-Gateway reverse proxy (for proto-defined HTTP endpoints)
// with optional hand-written routes (file uploads, WebSocket, SSE, health probes).
type HttpHandler struct {
	gatewayMux *runtime.ServeMux
	// srv service.Service — for non-proto endpoints
}

// New creates the HTTP handler and registers the gRPC-Gateway reverse proxy.
// grpcAddr is the address of the local gRPC server (e.g. "localhost:9090").
// All RPC methods annotated with google.api.http in the proto file are
// automatically exposed as REST/JSON endpoints through the gateway.
func New(ctx context.Context, grpcAddr string /*, srv Service*/) (*HttpHandler, error) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pbentity.RegisterServiceNameServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		return nil, err
	}

	return &HttpHandler{
		gatewayMux: mux,
		// srv: srv,
	}, nil
}

// Handler returns the http.Handler to use in app/run.go:
//
//	h, err := handlerhttp.New(ctx, cfg.GRPC.Addr)
//	httpServer := &http.Server{Addr: ":8080", Handler: h.Handler()}
func (h *HttpHandler) Handler() http.Handler {
	return h.gatewayMux
}
