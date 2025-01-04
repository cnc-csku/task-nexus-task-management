package grpcclient

type GrpcClientConfig struct {
	TaskManagementService grpcClientConfig `mapstructure:"taskManagementService"`
	NotificationService   grpcClientConfig `mapstructure:"notificationService"`
}

type grpcClientConfig struct {
	Name           string `mapstructure:"name"`
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	MaxSendMsgSize int    `mapstructure:"maxSendMsgSize"`
	MaxRecvMsgSize int    `mapstructure:"maxRecvMsgSize"`
}
