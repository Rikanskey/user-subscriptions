package logger

import (
	"bytes"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type StructuredLogger struct {
	Logger *logrus.Logger
}

func NewStructuredLogger(logger *logrus.Logger) func(handler http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

type StructuredLoggerEntry struct {
	Logger logrus.FieldLogger
}

const responseRounding = 100

func (s StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	s.Logger = s.Logger.WithFields(logrus.Fields{
		"resp_status":      status,
		"res_bytes_length": bytes,
		"resp_elapsed":     elapsed.Round(time.Millisecond / responseRounding).String(),
	})
	s.Logger.Info("Request completed")
}

func (s StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	s.Logger = s.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

func (s StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: logrus.NewEntry(s.Logger)}
	logFields := logrus.Fields{}

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["request_id"] = reqID
	}

	logFields["http_method"] = r.Method
	logFields["remote_address"] = r.RemoteAddr
	logFields["uri"] = r.RequestURI

	if requestBody := copyRequestBody(r); requestBody != "" {
		logFields["request_body"] = requestBody
	}

	entry.Logger = entry.Logger.WithFields(logFields)

	entry.Logger.Info("Request started")

	return entry
}

func copyRequestBody(r *http.Request) string {
	body, _ := io.ReadAll(r.Body)
	_ = r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	return string(body)
}

func GetLogEntry(r *http.Request) logrus.FieldLogger {
	entry, ok := middleware.GetLogEntry(r).(*StructuredLoggerEntry)
	if !ok {
		panic("LogEntry isn't *StructuredLoggerEntry")
	}

	return entry.Logger
}
