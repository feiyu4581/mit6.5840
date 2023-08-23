package option

import "time"

type Option struct {
	CoordinateAddress   string `mapstructure:"coordinate_address"`
	CurrentAddress      string `mapstructure:"current_address"`
	GrpcPort            int    `mapstructure:"grpc_port"`
	HeartBeatIntervalMs int    `mapstructure:"heart_beat_interval_ms"`
	HeartBeatIntervalS  int    `mapstructure:"heart_beat_interval_s"`
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
