package grpc_server

import (
	"context"
	"google.golang.org/grpc"
	"mapreduce/internal/map_server/service"
	"mapreduce/internal/pkg/rpc/map_server"
)

type Server struct {
	map_server.UnimplementedMapServer
}

func (s *Server) NewTask(ctx context.Context, taskRequest *map_server.Task) (*map_server.TaskStatus, error) {
	task := service.NewTask(taskRequest.TaskId, taskRequest.Name, taskRequest.Filename, taskRequest.Params)

	if err := service.GetWorkerManager().AddTask(task); err != nil {
		return nil, err
	}

	return &map_server.TaskStatus{
		TaskId: taskRequest.TaskId,
		Status: map_server.TaskStatus_Init,
	}, nil
}

func InitGrpcServer(server *grpc.Server) *grpc.Server {
	map_server.RegisterMapServer(server, &Server{})
	return server
}
