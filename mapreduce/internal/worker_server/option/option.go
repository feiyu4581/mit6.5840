package option

import (
	"mapreduce/internal/pkg/model"
	"time"
)

type Option struct {
	CoordinateAddress   string   `mapstructure:"coordinate_address"`
	CurrentAddress      string   `mapstructure:"current_address"`
	GrpcPort            int      `mapstructure:"grpc_port"`
	HeartBeatIntervalMs int      `mapstructure:"heart_beat_interval_ms"`
	HeartBeatIntervalS  int      `mapstructure:"heart_beat_interval_s"`
	Functions           []string `mapstructure:"functions"`
	FunctionRoute       string   `mapstructure:"function_route"`
	Mode                string   `mapstructure:"mode"`
}

func (option *Option) GetHeartBeatInterval() time.Duration {
	if option.HeartBeatIntervalMs > 0 {
		return time.Duration(option.HeartBeatIntervalMs) * time.Millisecond
	}
	if option.HeartBeatIntervalS > 0 {
		return time.Duration(option.HeartBeatIntervalS) * time.Second
	}
	return 5 * time.Second
}

func (option *Option) GetMode() model.WorkerMode {
	switch option.Mode {
	case "map":
		return model.MapMode
	case "reduce":
		return model.ReduceMode
	default:
		return model.MapMode
	}
}
