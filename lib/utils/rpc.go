package utils

import (
	"fmt"
	"net"
	"net/http"

	"golang.org/x/net/context"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"google.golang.org/grpc"
)

type ReverseProxyFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error

func PortToString(port int) string {
	return fmt.Sprintf(":%d", port)
}

func ServeGRPC(port int, server *grpc.Server) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		glog.Errorf("Error creating listener: %v", err)
		return err
	}
	glog.Infof("Serving GRPC on port %d", port)
	return server.Serve(lis)
}

func GrpcReverseProxy(ctx context.Context, fn ReverseProxyFunc, external, internal int, opts ...grpc.DialOption) error {
	jsonPbOpts := &runtime.JSONPb{OrigName: true, EmitDefaults: true}
	marshallingOpts := runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonPbOpts)
	mux := runtime.NewServeMux(marshallingOpts)
	if len(opts) == 0 {
		opts = append(opts, grpc.WithInsecure())
	}
	err := fn(ctx, mux, PortToString(internal), opts)
	if err != nil {
		return err
	}
	glog.Infof("Serving GRPC reverse-proxy on port %d", external)
	return http.ListenAndServe(PortToString(external), mux)
}
