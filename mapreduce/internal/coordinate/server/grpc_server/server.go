package grpc_server

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"mapreduce/internal/coordinate/model"
	"mapreduce/internal/coordinate/service/coordinator"
	"mapreduce/internal/pkg/log"
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

func (s *Server) Register(_ context.Context, info *coordinate.ClientInfo) (*coordinate.Response, error) {
	if err := coordinator.GetService().NewWorker(info.Name, info.Address, info.Mode); err != nil {
		return NewErrorResponse(err)
	}

	return NewSuccessResponse()
}

func (s *Server) Heartbeat(_ context.Context, request *coordinate.HeartbeatRequest) (*coordinate.Response, error) {
	wk, err := coordinator.GetService().GetWorker(request.Name)
	if err != nil {
		return NewErrorResponse(err)
	}

	switch request.Status {
	case coordinate.WorkerStatus_RunningStatus:
		wk.SetRunningStatus()
	case coordinate.WorkerStatus_OfflineStatus:
		wk.SetOfflineStatus()
	}

	wk.UpdateHeartbeat()

	log.Info("接收到 heatbeat 消息：%v", request)
	if request.CurrentTask != "" {
		task := coordinator.GetService().GetTask(request.CurrentTask)
		if task == nil {
			log.Error("Task <%s> does not exists", request.CurrentTask)
		} else {
			switch request.TaskStatus {
			case coordinate.TaskStatus_TaskRunningStatus:
				log.Info("Task <%s> is running", request.CurrentTask)
			case coordinate.TaskStatus_TaskFinishedStatus:
				if task.Status == model.RunningMapStatus {
					log.Info("赋值 map 结果：%s", request.Filenames)
					task.Status = model.WaitingReduceStatus
					task.MapWorkers[request.Name].FileNames = request.Filenames
				} else if task.Status == model.RunningReduceStatus {
					log.Info("赋值 reduce 结果：%s", request.Filenames)
					task.Status = model.FinishedStatus
					task.ReduceWorkers[request.Name].FileNames = request.Filenames
				}
			case coordinate.TaskStatus_TaskFailedStatus:
				task.Status = model.FailedStatus
				task.Err = errors.New(request.Message)
			}
			task.SingleC <- struct{}{}
		}
	}

	return NewSuccessResponse()
}

func InitGrpcServer(server *grpc.Server) *grpc.Server {
	coordinate.RegisterCoordinateServer(server, &Server{})

	return server
}
