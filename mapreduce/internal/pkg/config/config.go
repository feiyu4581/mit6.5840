package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func UnmarshalConfig(configAddress string, filename string, configObj interface{}) {
	cfg := viper.New()
	cfg.SetConfigName(filename)
	cfg.SetConfigType("yaml")
	cfg.AddConfigPath(configAddress)

	if err := cfg.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("config read error: %s", err.Error()))
	}

	if err := cfg.Unmarshal(configObj); err != nil {
		panic(fmt.Sprintf("config unmarshal error: %s", err.Error()))
	}
}
