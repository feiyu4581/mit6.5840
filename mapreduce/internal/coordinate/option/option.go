package option

type Option struct {
	Port     int `mapstructure:"port"`
	GrpcPort int `mapstructure:"grpc_port"`
}
