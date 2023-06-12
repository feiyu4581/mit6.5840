package grpc_server

import (
	"google.golang.org/grpc"
	"mapreduce/internal/coordinate/server/grpc_server/pb"
)

func InitGrpcServer(server *grpc.Server) *grpc.Server {
	pb.RegisterMapServiceServer(server, MapService{})

	return server
}
