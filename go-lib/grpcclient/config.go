package grpcclient

type GrpcClientConfig struct {
	TaskManagementService grpcClientConfig `envPrefix:"TASK_MANAGEMENT_SERVICE_"`
	NotificationService   grpcClientConfig `envPrefix:"NOTIFICATION_SERVICE_"`
}

type grpcClientConfig struct {
	Name           string `env:"NAME"`
	Host           string `env:"HOST"`
	Port           int    `env:"PORT"`
	MaxSendMsgSize int    `env:"MAX_SEND_MSG_SIZE"`
	MaxRecvMsgSize int    `env:"MAX_RECV_MSG_SIZE"`
}
