package server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	api_pb "github.com/gofunct/runtime/testing/api"
)

// RegisterWithServer implements runtime.Server.RegisterWithServer.
func (s *libraryServiceServerImpl) RegisterWithServer(grpcSvr *grpc.Server) {
	api_pb.RegisterLibraryServiceServer(grpcSvr, s)
}

// RegisterWithHandler implements runtime.Server.RegisterWithHandler.
func (s *libraryServiceServerImpl) RegisterWithHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return api_pb.RegisterLibraryServiceHandler(ctx, mux, conn)
}
