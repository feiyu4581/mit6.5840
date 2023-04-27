package server

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"mapreduce/internal/coordinate/option"
	"mapreduce/internal/coordinate/server/grpc_server"
	"mapreduce/internal/coordinate/server/http_server"
	"mapreduce/internal/pkg/log"
	"net"
	"net/http"
)

type HttpServer struct {
	Option     *option.Option
	Server     *http.Server
	GrpcServer *grpc.Server
}

func NewServer(ctx context.Context, option *option.Option) (*HttpServer, error) {
	server := &HttpServer{
		Option: option,
	}

	if err := server.Init(); err != nil {
		return nil, errors.Wrap(err, "server init error")
	}

	return server, nil
}

func (server *HttpServer) GetOption() *option.Option {
	return server.Option
}

func (server *HttpServer) Init() error {
	log.InitLogrus()
	server.Server = http_server.NewHttpServer(server.GetOption().Port)
	server.GrpcServer = grpc_server.InitGrpcServer(grpc.NewServer())
	return nil
}

func (server *HttpServer) Start() {
	go func() {
		log.Info("http server: %d", server.GetOption().Port)
		if err := server.Server.ListenAndServe(); err != nil {
			log.Fatal("http server error: %s", err.Error())
		}
	}()

	go func() {
		log.Info("grpc server: %d", server.GetOption().GrpcPort)
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", server.GetOption().GrpcPort))
		if err != nil {
			log.Fatal("grpc failed error: %s", err.Error())
		}

		if err = server.GrpcServer.Serve(lis); err != nil {
			log.Fatal("grpc server error: %s", err.Error())
		}
	}()
}
