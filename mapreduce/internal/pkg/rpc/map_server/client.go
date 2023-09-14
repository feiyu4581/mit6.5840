package map_server

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

var (
	clientMaps = make(map[string]*ServerClient, 0)
	clientLock sync.Mutex
)

type ServerClient struct {
	address string
	client  MapClient
}

func GetServerClient(address string) (*ServerClient, error) {
	clientLock.Lock()
	defer clientLock.Unlock()

	if client, ok := clientMaps[address]; ok {
		return client, nil
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("grpc dial %s error", address))
	}

	clientMaps[address] = &ServerClient{
		address: address,
		client:  NewMapClient(conn),
	}

	return clientMaps[address], nil
}

func (client *ServerClient) NewTask(ctx context.Context, name string, filenames []string, function string, splitNums, taskId, taskIndex int64) error {
	response, err := client.client.NewTask(ctx, &Task{
		Name:      name,
		Filenames: filenames,
		Function:  function,
		SplitNums: splitNums,
		TaskId:    taskId,
		TaskIndex: taskIndex,
	})

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("grpc new task error"))
	}

	if !response.Success {
		return errors.New(response.Message)
	}

	return nil
}
