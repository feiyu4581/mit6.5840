package cmd

import (
	"context"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"mapreduce/internal/map_server/option"
	MapServer "mapreduce/internal/map_server/server"
	"mapreduce/internal/pkg/config"
	"mapreduce/internal/pkg/log"
	"mapreduce/internal/pkg/rpc/coordinate"
	"os"
	"os/signal"
	"syscall"
)

var mapOption option.Option
var workerName string

func initMapOption() {
	workerName = uuid.New().String()
	config.UnmarshalConfig("./config/test/", "map_server.yaml", &mapOption)
}

var mapCmd = &cobra.Command{
	Use:   "map [command]",
	Short: "Start a map server",
	Run: func(cmd *cobra.Command, args []string) {
		initMapOption()

		ctx, cancel := context.WithCancel(context.Background())
		server, err := MapServer.NewServer(ctx, workerName, &mapOption)
		if err != nil {
			log.Fatal("Server init error: %s", err.Error())
		}

		err = coordinate.GetServerClient(mapOption.CoordinateAddress).Register(ctx, workerName, mapOption.CurrentAddress)
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
