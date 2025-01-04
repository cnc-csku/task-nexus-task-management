package config

import core_grpcclient "github.com/cnc-csku/task-nexus/go-lib/grpcclient"

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
