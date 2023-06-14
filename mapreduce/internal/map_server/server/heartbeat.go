package server

import (
	"mapreduce/internal/pkg/log"
	"mapreduce/internal/pkg/rpc/coordinate"
	"time"
)

func (server *MapServer) LoopForHeartBeat() {
	ticker := time.NewTicker(server.GetOption().GetHeartBeatInterval())
	client := coordinate.GetServerClient(server.GetOption().CoordinateAddress)
	for {
		select {
		case <-ticker.C:
			if err := client.HeartbeatRunning(server.Ctx, server.Name); err != nil {
				log.Error("heartbeat error: %s", err.Error())
			}
		case <-server.Ctx.Done():
			if err := client.HeartbeatOffline(server.Ctx, server.Name); err != nil {
				log.Error("heartbeat error: %s", err.Error())
			}
		}
		log.Info("done")
	}
}
