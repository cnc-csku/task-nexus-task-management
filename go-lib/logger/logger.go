package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func NewLogrusLogger() *logrus.Logger {
	logger := &logrus.Logger{
		Out:   os.Stdout,
		Level: logrus.DebugLevel,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			PrettyPrint:     true,
		},
	}

	return logger
}
