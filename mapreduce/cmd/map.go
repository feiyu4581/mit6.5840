package cmd

import (
	"context"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"mapreduce/internal/pkg/config"
	"mapreduce/internal/pkg/log"
	"mapreduce/internal/pkg/model"
	"mapreduce/internal/pkg/rpc/coordinate"
	"mapreduce/internal/worker_server/option"
	MapServer "mapreduce/internal/worker_server/server"
	"os"
	"os/signal"
	"syscall"
)

var mapOption option.Option
var workerName string

func initMapOption(args []string) {
	workerName = uuid.New().String()
	filename := "map_server.yaml"
	if len(args) > 0 {
		filename = args[0]
	}
	config.UnmarshalConfig("./config/test/", filename, &mapOption)
}

var mapCmd = &cobra.Command{
	Use:   "map [command]",
	Short: "Start a map server",
	Run: func(cmd *cobra.Command, args []string) {
		initMapOption(args)

		ctx, cancel := context.WithCancel(context.Background())
		server, err := MapServer.NewServer(ctx, workerName, &mapOption)
		if err != nil {
			log.Fatal("Server init error: %s", err.Error())
		}

		mode := coordinate.ClientMode_MapMode
		if mapOption.GetMode() == model.ReduceMode {
			mode = coordinate.ClientMode_ReduceMode
		}

		err = coordinate.GetServerClient(mapOption.CoordinateAddress).Register(ctx, workerName, mapOption.CurrentAddress, mode)
		if err != nil {
			log.Fatal("Register coordinate error: %s", err.Error())
		}

		server.Start()

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, os.Kill)
		<-signalChan

		cancel()
		server.GrpcServer.GracefulStop()
	},
}
