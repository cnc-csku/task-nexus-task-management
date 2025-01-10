package logging

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type responseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *responseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func EchoLoggingMiddleware(logger *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			start := time.Now()
			res := c.Response()

			// Create a buffer to capture the response body
			var resBody bytes.Buffer
			mw := io.MultiWriter(res.Writer, &resBody)
			writer := &responseWriter{Writer: mw, ResponseWriter: res.Writer}
			res.Writer = writer

			err := next(c)
			stop := time.Now()

			latency := stop.Sub(start)
			statusCode := res.Status

			fields := logrus.Fields{
				"type":       "httpserver",
				"method":     req.Method,
				"uri":        req.RequestURI,
				"remote_ip":  c.RealIP(),
				"status":     statusCode,
				"latency":    latency.String(),
				"user_agent": req.UserAgent(),
			}

			var requestObj map[string]interface{}
			if err := c.Bind(&requestObj); err != nil {
				fields["request"] = "error parsing request"
			} else {
				fields["request"] = requestObj
			}

			var level logrus.Level
			if err != nil {
				level = logrus.ErrorLevel
				fields["error"] = err.Error()
			} else {
				level = logrus.InfoLevel
				fields["response_body"] = resBody.String()
			}

			logger.WithFields(fields).Log(level, "request completed")

			return err
		}
	}
}
