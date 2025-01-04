package config

import (
	"log"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	// Service
	ServiceName string `mapstructure:"SERVICE_NAME"`

	// Rest Server
	RestPort string `mapstructure:"REST_PORT"`

	// Database
	MongoURI string `mapstructure:"MONGO_URI"`

	// gRPC Server
	GrpcPort          string `mapstructure:"GRPC_PORT"`
	GrpcMaxSendSize   int    `mapstructure:"GRPC_MAX_SEND_SIZE"`
	GrpcMaxRecvSize   int    `mapstructure:"GRPC_MAX_RECV_SIZE"`
	GrpcUseReflection bool   `mapstructure:"GRPC_USE_REFLECTION"`
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

func NewConfig() *Config {
	config := Config{}

	// Set the .env file and read environment variables
	viper.SetConfigFile(".env")

	// Attempt to read the .env file
	if err := viper.ReadInConfig(); err != nil {
		log.Println("⚠️  .env file not found or cannot be read, using environment variables")

		// Bind environment variables
		envs := getMapstructureTags(config)
		for _, env := range envs {
			viper.MustBindEnv(env)
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Println("❌ Unable to decode into struct:", err)
	}

	return &config
}
