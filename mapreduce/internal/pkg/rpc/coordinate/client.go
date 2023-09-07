package coordinate

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mapreduce/internal/coordinate/model"
	mapModel "mapreduce/internal/worker_server/model"
	"mapreduce/internal/worker_server/service"
	"sync"
)

var (
	client         *ServerClient
	coordinateOnce sync.Once
)

type ServerClient struct {
	client CoordinateClient
}

func GetServerClient(address string) *ServerClient {
	coordinateOnce.Do(func() {
		conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}

		client = &ServerClient{
			client: NewCoordinateClient(conn),
		}
	})

	return client
}

func (client *ServerClient) Register(ctx context.Context, name string, address string, mode ClientMode) error {
	response, err := client.client.Register(ctx, &ClientInfo{
		Name:    name,
		Address: address,
		Mode:    mode,
	})

	if err != nil {
		return errors.Wrap(err, "register error")
	}

	if response.Code != 0 {
		return errors.New(response.Message)
	}

	return nil
}

func (client *ServerClient) HeartbeatRunning(ctx context.Context, wk *model.Worker) error {
	return client.Heartbeat(ctx, wk, WorkerStatus_RunningStatus)
}

func (client *ServerClient) HeartbeatOffline(ctx context.Context, wk *model.Worker) error {
	return client.Heartbeat(ctx, wk, WorkerStatus_OfflineStatus)
}

func (client *ServerClient) Heartbeat(ctx context.Context, wk *model.Worker, status WorkerStatus) error {
	workerManager := service.GetWorkerManager()
	request := &HeartbeatRequest{
		Name:   wk.Name,
		Status: status,
	}

	if workerManager != nil && workerManager.CurrentTask != nil {
		request.CurrentTask = workerManager.CurrentTask.Name
		switch workerManager.CurrentTask.Status {
		case mapModel.InitStatus:
		case mapModel.RunningStatus:
			request.TaskStatus = TaskStatus_TaskRunningStatus
		case mapModel.DoneStatus:
			request.TaskStatus = TaskStatus_TaskFinishedStatus
			request.Filenames = workerManager.CurrentTask.ResFileNames
		case mapModel.FailedStatus:
			request.TaskStatus = TaskStatus_TaskFailedStatus
			request.Message = workerManager.CurrentTask.Err.Error()
		}
	}

	response, err := client.client.Heartbeat(ctx, request)

	if err != nil {
		return errors.Wrap(err, "heartbeat error")
	}

	if response.Code != 0 {
		return errors.New(response.Message)
	}

	return nil
}
