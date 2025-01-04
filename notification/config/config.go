package config

import (
	"log"

	"github.com/cnc-csku/task-nexus/go-lib/config"
	core_grpcclient "github.com/cnc-csku/task-nexus/go-lib/grpcclient"
	"github.com/spf13/viper"
)

type Config struct {
	config.Config `mapstructure:",squash"` // squash the nested struct into the parent struct
}

func NewConfig() *Config {
	config := &Config{}

	viper.SetConfigName("config.yaml")
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("Error reading config file", err)
	}

	if err := viper.Unmarshal(config); err != nil {
		log.Fatalln("Unable to decode into struct", err)
	}

	return config
}

func ProvideGrpcClientConfig(config *Config) core_grpcclient.GrpcClientConfig {
	return config.GrpcClient
}
