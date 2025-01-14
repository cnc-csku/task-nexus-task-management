package logging

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cnc-csku/task-nexus/go-lib/utils/errutils"
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
	logger.SetFormatter(&CustomFormatter{})

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			start := time.Now()
			res := c.Response()

			// Save the request body
			var reqBody []byte
			if req.Body != nil {
				reqBody, _ = io.ReadAll(req.Body)
				req.Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Restore the request body
			}
			defer req.Body.Close()

			// Create a buffer to capture the response body
			var resBody bytes.Buffer
			mw := io.MultiWriter(res.Writer, &resBody)
			writer := &responseWriter{Writer: mw, ResponseWriter: res.Writer}
			res.Writer = writer

			err := next(c)

			// Check if the error is of type *echo.HTTPError and unwrap the original error
			if httpErr, ok := err.(*echo.HTTPError); ok {
				if httpErr.Internal != nil {
					err = httpErr.Internal
				}
			}

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

			// Unmarshal request body to remove escape characters
			var requestObj map[string]interface{}
			if err := json.Unmarshal(reqBody, &requestObj); err != nil {
				fields["request_body"] = "error parsing request"
			} else {
				fields["request_body"] = requestObj
			}

			var level logrus.Level
			if err != nil {
				level = logrus.ErrorLevel

				serviceError, ok := err.(*errutils.Error)
				if ok && serviceError != nil {
					stackTraceErr := errutils.GetStackField(serviceError.StackErr)
					stackTraceParts := strings.Split(stackTraceErr.Stack, "\n\t")
					if len(stackTraceParts) > 10 {
						stackTraceErr.Stack = strings.Join(stackTraceParts[:10], "\n\t")
					}
					fields["stack_trace"] = logrus.Fields{
						"type":  stackTraceErr.Type,
						"stack": stackTraceErr.Stack,
					}
					fields["error"] = logrus.Fields{
						"status":        serviceError.Status.String(),
						"message":       serviceError.Message,
						"debug_message": serviceError.DebugMessage,
					}
				} else {
					fields["error"] = logrus.Fields{
						"message": err.Error(),
					}
				}
			} else {
				level = logrus.InfoLevel
				var responseBody map[string]interface{}
				if err := json.Unmarshal(resBody.Bytes(), &responseBody); err == nil {
					fields["response_body"] = responseBody
				} else {
					fields["response_body"] = resBody.String()
				}
			}

			logger.WithFields(fields).Log(level, "Logging incoming request")

			return err
		}
	}
}
