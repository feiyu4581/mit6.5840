package grpc_server

import (
	"context"
	"mapreduce/internal/coordinate/server/grpc_server/pb"
	"mapreduce/internal/coordinate/service/map_worker"
)

type MapService struct {
	pb.UnimplementedMapServiceServer
}

func NewResponse(code int64, message string) *pb.Response {
	return &pb.Response{
		Code:    code,
		Message: message,
	}
}

func NewSuccessResponse() (*pb.Response, error) {
	return NewResponse(0, "ok"), nil
}

func NewErrorResponse(err error) (*pb.Response, error) {
	return NewResponse(-1, err.Error()), nil
}

func (m MapService) Register(ctx context.Context, info *pb.ClientInfo) (*pb.Response, error) {
	_, err := map_worker.NewWorker(info.Name, info.Address)
	if err != nil {
		return NewErrorResponse(err)
	}

	return NewSuccessResponse()
}

func (m MapService) Heartbeat(ctx context.Context, request *pb.HeartbeatRequest) (*pb.Response, error) {
	worker, err := map_worker.GetWorker(request.Name)
	if err != nil {
		return NewErrorResponse(err)
	}

	switch request.Status {
	case pb.WorkerStatus_RunningStatus:
		worker.SetRunningStatus()
	case pb.WorkerStatus_OfflineStatus:
		worker.SetOfflineStatus()
	}

	worker.UpdateHeartbeat()
	return NewSuccessResponse()
}
