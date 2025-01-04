package config

import (
	"log"
	"reflect"

	core_grpcclient "github.com/cnc-csku/task-nexus/go-lib/grpcclient"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceName string                           `mapstructure:"serviceName"`
	RestServer  RestServerConfig                 `mapstructure:"restServer"`
	MongoDB     MongoDBConfig                    `mapstructure:"mongoDB"`
	GrpcServer  GrpcServerConfig                 `mapstructure:"grpcServer"`
	GrpcClient  core_grpcclient.GrpcClientConfig `mapstructure:"grpcClient"`
}

type RestServerConfig struct {
	Port string `mapstructure:"port"`
}

type MongoDBConfig struct {
	URI      string `mapstructure:"uri"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type GrpcServerConfig struct {
	Port           string `mapstructure:"port"`
	MaxSendMsgSize int    `mapstructure:"maxSendMsgSize"`
	MaxRecvMsgSize int    `mapstructure:"maxRecvMsgSize"`
	UseReflection  bool   `mapstructure:"useReflection"`
}

func NewConfig() *Config {
	config := &Config{}

	viper.SetConfigName("config.yaml")
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("⚠️  .env file not found or cannot be read, using environment variables")

		// Bind environment variables
		envs := getMapstructureTags(config)
		for _, env := range envs {
			viper.MustBindEnv(env)
		}
	}

	if err := viper.Unmarshal(config); err != nil {
		log.Fatalln("Unable to decode into struct", err)
	}

	return config
}

func getMapstructureTags(v interface{}) []string {
	typ := reflect.TypeOf(v)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	var tags []string
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if tag, ok := field.Tag.Lookup("mapstructure"); ok {
			tags = append(tags, tag)
		}
	}
	return tags
}

func ProvideGrpcClientConfig(config *Config) core_grpcclient.GrpcClientConfig {
	return config.GrpcClient
}
