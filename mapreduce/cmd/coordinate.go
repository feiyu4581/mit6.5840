package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"mapreduce/internal/coordinate/option"
	CoordinateServer "mapreduce/internal/coordinate/server"
	"mapreduce/internal/pkg/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var coordinateOption option.Option

var coordinateCmd = &cobra.Command{
	Use:   "coordinate [command]",
	Short: "Start a coordinate server",
	Run: func(cmd *cobra.Command, args []string) {
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

func init() {
	coordinateCmd.PersistentFlags().IntVarP(
		&coordinateOption.Port, "port", "p", 5432, "Coordinate server port")

	coordinateCmd.PersistentFlags().IntVarP(
		&coordinateOption.GrpcPort, "grpc_port", "g", 5433, "Coordinate grpc server port")
}
