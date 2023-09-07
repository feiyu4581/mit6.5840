package server

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"mapreduce/internal/coordinate/model"
	"mapreduce/internal/pkg/log"
	functionModel "mapreduce/internal/pkg/model"
	"mapreduce/internal/worker_server/option"
	"mapreduce/internal/worker_server/server/grpc_server"
	"mapreduce/internal/worker_server/service"
	"net"
	"plugin"
)

type MapServer struct {
	Option     *option.Option
	Ctx        context.Context
	GrpcServer *grpc.Server
	Worker     *model.Worker
}

func NewServer(ctx context.Context, workerName string, option *option.Option) (*MapServer, error) {
	wk := model.NewWorker(workerName, fmt.Sprintf(":%d", option.GrpcPort))
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
	server.RegisterFunction()

	service.GetWorkerManager().SetWorkMode(server.Option.GetMode())

	return nil
}

func (server *MapServer) RegisterFunction() {
	for _, functionName := range server.Option.Functions {
		functionFilePath := fmt.Sprintf("%s/%s.so", server.Option.FunctionRoute, functionName)
		p, err := plugin.Open(functionFilePath)
		if err != nil {
			panic(fmt.Sprintf("open plugin error: %s", err.Error()))
		}

		xMapF, err := p.Lookup("Map")
		if err != nil {
			panic(fmt.Sprintf("lookup map function error: %s", err.Error()))
		}

		mapf := xMapF.(func(string, string) []functionModel.KeyValue)

		xReduceF, err := p.Lookup("Reduce")
		if err != nil {
			panic(fmt.Sprintf("lookup reduce function error: %s", err.Error()))
		}

		log.Info("register function <%s> success", functionName)
		reduceF := xReduceF.(func(string, []string) string)
		service.GetWorkerManager().Register(&functionModel.Function{
			Name:   functionName,
			Map:    mapf,
			Reduce: reduceF,
		})
	}
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
