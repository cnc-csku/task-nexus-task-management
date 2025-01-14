package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b bytes.Buffer

	// Format time
	timestamp := entry.Time.Format(time.RFC3339Nano)
	b.WriteString(fmt.Sprintf("%s\t", timestamp))

	// Format level with color
	level := strings.ToUpper(entry.Level.String())
	coloredLevel := f.colorizeLevel(level, entry.Level)
	b.WriteString(fmt.Sprintf("%s\t", coloredLevel))

	// Format message
	b.WriteString(fmt.Sprintf("%s\t", entry.Message))

	// Format fields as JSON
	fields, err := json.Marshal(entry.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON: %w", err)
	}
	b.WriteString(string(fields))

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *CustomFormatter) colorizeLevel(level string, logLevel logrus.Level) string {
	switch logLevel {
	case logrus.DebugLevel:
		return fmt.Sprintf("\033[36m%s\033[0m", level) // Cyan
	case logrus.InfoLevel:
		return fmt.Sprintf("\033[32m%s\033[0m", level) // Green
	case logrus.WarnLevel:
		return fmt.Sprintf("\033[33m%s\033[0m", level) // Yellow
	case logrus.ErrorLevel:
		return fmt.Sprintf("\033[31m%s\033[0m", level) // Red
	case logrus.FatalLevel:
		return fmt.Sprintf("\033[35m%s\033[0m", level) // Magenta
	case logrus.PanicLevel:
		return fmt.Sprintf("\033[41m%s\033[0m", level) // Red background
	default:
		return level
	}
}
