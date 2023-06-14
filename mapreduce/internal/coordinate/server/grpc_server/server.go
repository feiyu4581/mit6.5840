package grpc_server

import (
	"context"
	"google.golang.org/grpc"
	"mapreduce/internal/coordinate/service/map_worker"
	"mapreduce/internal/pkg/rpc/coordinate"
)

type Server struct {
	coordinate.UnimplementedCoordinateServer
}

func NewResponse(code int64, message string) *coordinate.Response {
	return &coordinate.Response{
		Code:    code,
		Message: message,
	}
}

func NewSuccessResponse() (*coordinate.Response, error) {
	return NewResponse(0, "ok"), nil
}

func NewErrorResponse(err error) (*coordinate.Response, error) {
	return NewResponse(-1, err.Error()), nil
}

func (s *Server) Register(ctx context.Context, info *coordinate.ClientInfo) (*coordinate.Response, error) {
	_, err := map_worker.NewWorker(info.Name, info.Address)
	if err != nil {
		return NewErrorResponse(err)
	}

	return NewSuccessResponse()
}

func (s *Server) Heartbeat(ctx context.Context, request *coordinate.HeartbeatRequest) (*coordinate.Response, error) {
	worker, err := map_worker.GetWorker(request.Name)
	if err != nil {
		return NewErrorResponse(err)
	}

	switch request.Status {
	case coordinate.WorkerStatus_RunningStatus:
		worker.SetRunningStatus()
	case coordinate.WorkerStatus_OfflineStatus:
		worker.SetOfflineStatus()
	}

	worker.UpdateHeartbeat()
	return NewSuccessResponse()
}

func InitGrpcServer(server *grpc.Server) *grpc.Server {
	coordinate.RegisterCoordinateServer(server, &Server{})

	return server
}
