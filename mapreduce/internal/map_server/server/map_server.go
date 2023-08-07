package server

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"mapreduce/internal/coordinate/service/worker"
	"mapreduce/internal/map_server/option"
	"mapreduce/internal/map_server/server/grpc_server"
	"mapreduce/internal/pkg/log"
	"net"
)

type MapServer struct {
	Option     *option.Option
	Ctx        context.Context
	GrpcServer *grpc.Server
	Worker     *worker.Worker
}

func NewServer(ctx context.Context, workerName string, option *option.Option) (*MapServer, error) {
	wk, err := worker.NewWorker(workerName, fmt.Sprintf(":%d", option.GrpcPort))
	if err != nil {
		return nil, errors.Wrap(err, "new worker error")
	}

	server := &MapServer{
		Option: option,
		Worker: wk,
		Ctx:    ctx,
	}

	if err := server.Init(); err != nil {
		return nil, errors.Wrap(err, "server init error")
	}

	return server, nil
}

func (server *MapServer) GetOption() *option.Option {
	return server.Option
}

func (server *MapServer) Init() error {
	log.InitLogrus()
	server.GrpcServer = grpc_server.InitGrpcServer(grpc.NewServer())
	return nil
}

func (server *MapServer) Start() {
	go func() {
		log.Info("grpc server: %d", server.GetOption().GrpcPort)
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", server.GetOption().GrpcPort))
		if err != nil {
			log.Fatal("grpc failed error: %s", err.Error())
		}

		if err = server.GrpcServer.Serve(lis); err != nil {
			log.Fatal("grpc failed error: %s", err.Error())
		}
	}()

	go server.LoopForHeartBeat()
}
