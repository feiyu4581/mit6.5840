package option

import "time"

type Option struct {
	CoordinateAddress   string
	CurrentAddress      string
	GrpcPort            int
	HeartBeatIntervalMs int
	HeartBeatIntervalS  int
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
