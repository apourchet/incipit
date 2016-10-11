package utils

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

func ServeGRPC(port int, server *grpc.Server) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	return server.Serve(lis)
}
