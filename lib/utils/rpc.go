package utils

import (
	"fmt"
	"net"

	"github.com/golang/glog"

	"google.golang.org/grpc"
)

func ServeGRPC(port int, server *grpc.Server) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		glog.Errorf("Error creating listener: %v", err)
		return err
	}
	glog.Infof("Serving GRPC on port %d", port)
	return server.Serve(lis)
}
