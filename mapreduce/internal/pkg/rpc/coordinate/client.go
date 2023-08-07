package coordinate

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mapreduce/internal/coordinate/service/worker"
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

func (client *ServerClient) Register(ctx context.Context, name string, address string) error {
	response, err := client.client.Register(ctx, &ClientInfo{
		Name:    name,
		Address: address,
	})

	if err != nil {
		return errors.Wrap(err, "register error")
	}

	if response.Code != 0 {
		return errors.New(response.Message)
	}

	return nil
}

func (client *ServerClient) HeartbeatRunning(ctx context.Context, wk *worker.Worker) error {
	return client.Heartbeat(ctx, wk, WorkerStatus_RunningStatus)
}

func (client *ServerClient) HeartbeatOffline(ctx context.Context, wk *worker.Worker) error {
	return client.Heartbeat(ctx, wk, WorkerStatus_OfflineStatus)
}

func (client *ServerClient) Heartbeat(ctx context.Context, wk *worker.Worker, status WorkerStatus) error {
	response, err := client.client.Heartbeat(ctx, &HeartbeatRequest{
		Name:   wk.Name,
		Status: status,
	})

	if err != nil {
		return errors.Wrap(err, "heartbeat error")
	}

	if response.Code != 0 {
		return errors.New(response.Message)
	}

	return nil
}
