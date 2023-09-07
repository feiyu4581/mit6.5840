package grpc_server

import (
	"context"
	"google.golang.org/grpc"
	"mapreduce/internal/pkg/rpc/map_server"
	"mapreduce/internal/worker_server/model"
	"mapreduce/internal/worker_server/service"
)

type Server struct {
	map_server.UnimplementedMapServer
}

func (s *Server) NewTask(ctx context.Context, taskRequest *map_server.Task) (*map_server.TaskResponse, error) {
	task := model.NewTask(taskRequest.TaskId, taskRequest.TaskIndex, taskRequest.Name, taskRequest.Filename, taskRequest.SplitNums)

	response := &map_server.TaskResponse{Success: true}
	err := service.GetWorkerManager().AddTask(task, taskRequest.Function)
	if err != nil {
		response.Success = false
		response.Message = err.Error()
	}

	return response, nil
}

func InitGrpcServer(server *grpc.Server) *grpc.Server {
	map_server.RegisterMapServer(server, &Server{})
	return server
}