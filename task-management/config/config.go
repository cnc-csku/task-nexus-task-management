package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	coreGrpcClient "github.com/cnc-csku/task-nexus/go-lib/grpcclient"
	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName  string                          `env:"SERVICE_NAME"`
	AllowOrigins []string                        `env:"ALLOW_ORIGINS" envSeparator:","`
	RestServer   RestServerConfig                `envPrefix:"REST_SERVER_"`
	MongoDB      MongoDBConfig                   `envPrefix:"MONGO_"`
	GrpcServer   GrpcServerConfig                `envPrefix:"GRPC_SERVER_"`
	GrpcClient   coreGrpcClient.GrpcClientConfig `envPrefix:"GRPC_CLIENT_"`
	OllamaClient OllamaClientConfig              `envPrefix:"OLLAMA_CLIENT_"`
	JWT          JWT                             `envPrefix:"JWT_"`
	LogFormat    string                          `env:"LOG_FORMAT"`
}

type RestServerConfig struct {
	Port string `env:"PORT"`
}

type MongoDBConfig struct {
	URI      string `env:"URI"`
	Port     string `env:"PORT"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	Database string `env:"DATABASE"`
}

type GrpcServerConfig struct {
	Port           string `env:"PORT"`
	MaxSendMsgSize int    `env:"MAX_SEND_MSG_SIZE"`
	MaxRecvMsgSize int    `env:"MAX_RECV_MSG_SIZE"`
	UseReflection  bool   `env:"USE_REFLECTION"`
}

type OllamaClientConfig struct {
	Endpoint      string `env:"ENDPOINT"`
	UseProxy      bool   `env:"USE_PROXY"`
	HttpProxyHost string `env:"HTTP_PROXY_HOST"`
	HttpProxyPort string `env:"HTTP_PROXY_PORT"`
}

type JWT struct {
	AccessTokenSecret  string `env:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret string `env:"REFRESH_TOKEN_SECRET"`
}

func NewConfig() *Config {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it. Falling back to system environment variables.")
	}

	config := &Config{}

	// Parse environment variables into the config struct
	if err := env.Parse(config); err != nil {
		log.Fatalln("Failed to parse environment variables into Config struct:", err)
	}

	return config
}

func ProvideGrpcClientConfig(config *Config) coreGrpcClient.GrpcClientConfig {
	return config.GrpcClient
}
