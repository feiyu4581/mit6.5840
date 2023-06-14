package cmd

import (
	"context"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"mapreduce/internal/map_server/option"
	MapServer "mapreduce/internal/map_server/server"
	"mapreduce/internal/pkg/log"
	"mapreduce/internal/pkg/rpc/coordinate"
	"os"
	"os/signal"
	"syscall"
)

var mapOption option.Option
var workerName string

var mapCmd = &cobra.Command{
	Use:   "map [command]",
	Short: "Start a map server",
	Run: func(cmd *cobra.Command, args []string) {
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

func init() {
	workerName = uuid.New().String()

	mapCmd.PersistentFlags().IntVarP(
		&mapOption.GrpcPort, "grpc_port", "g", 6433, "Map grpc server port")

	mapCmd.PersistentFlags().IntVarP(
		&mapOption.HeartBeatIntervalMs, "heartbeat_interval_ms", "i", 5000, "Map heartbeat interval(millisecond)")

	mapCmd.PersistentFlags().IntVarP(
		&mapOption.HeartBeatIntervalS, "heartbeat_interval_s", "s", 5000, "Map heartbeat interval(second)")

	mapCmd.PersistentFlags().StringVarP(
		&mapOption.CoordinateAddress, "coordinate_address", "c", "127.0.0.1:5433", "Coordinate address")

	mapCmd.PersistentFlags().StringVarP(
		&mapOption.CurrentAddress, "current_address", "m", "127.0.0.1:6433", "Current address")
}
