package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"mapreduce/internal/coordinate/option"
	CoordinateServer "mapreduce/internal/coordinate/server"
	"mapreduce/internal/pkg/config"
	"mapreduce/internal/pkg/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var coordinateOption option.Option

func initCoordinatorOption() {
	config.UnmarshalConfig("./config/test/", "coordinate.yaml", &coordinateOption)
}

var coordinateCmd = &cobra.Command{
	Use:   "coordinate [command]",
	Short: "Start a coordinate server",
	Run: func(cmd *cobra.Command, args []string) {
		initCoordinatorOption()

		ctx := context.Background()
		server, err := CoordinateServer.NewServer(ctx, &coordinateOption)
		if err != nil {
			log.Fatal("Server init error: %s", err.Error())
		}

		server.Start()

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, os.Kill)
		<-signalChan

		timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := server.Server.Shutdown(timeoutCtx); err != nil {
			log.Fatal("CoordinateServer Shutdown:", err)
		}

		server.GrpcServer.GracefulStop()
		log.Info("Server close success")
	},
}
