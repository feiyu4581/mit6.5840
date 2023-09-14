package server

import (
	"mapreduce/internal/pkg/log"
	"mapreduce/internal/pkg/rpc/coordinate"
	"mapreduce/internal/worker_server/service"
	"time"
)

func (server *MapServer) LoopForHeartBeat() {
	ticker := time.NewTicker(server.GetOption().GetHeartBeatInterval())
	client := coordinate.GetServerClient(server.GetOption().CoordinateAddress)
	workerManager := service.GetWorkerManager()
	for {
		select {
		case <-workerManager.SingleC:
			if err := client.HeartbeatRunning(server.Ctx, server.Worker); err != nil {
				log.Error("heartbeat error: %s", err.Error())
			}
		case <-ticker.C:
			if err := client.HeartbeatRunning(server.Ctx, server.Worker); err != nil {
				log.Error("heartbeat error: %s", err.Error())
			}
		case <-server.Ctx.Done():
			if err := client.HeartbeatOffline(server.Ctx, server.Worker); err != nil {
				log.Error("heartbeat error: %s", err.Error())
			}
		}
	}
}
