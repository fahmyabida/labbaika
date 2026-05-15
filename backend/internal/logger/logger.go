package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// Key type for context values
type contextKey string

const XTraceID contextKey = "x-trace-id"

type logger struct {
	*logrus.Logger
}

var (
	staticLogger logger
)

func InitStaticLogger(serviceName string) {
	log := logrus.New()

	log.SetFormatter(&logrus.JSONFormatter{
		// custom timestamp to mirror .toISOString() in Javascript
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:  "message",
			logrus.FieldKeyTime: "event_timestamp",
		},
	})

	log.SetOutput(os.Stdout)

	staticLogger = logger{log}
}

func Info(ctx context.Context, messageMap map[string]interface{}, message string) {
	fields := buildFields(ctx, messageMap)
	staticLogger.WithFields(fields).Info(message)
}

func Infof(ctx context.Context, messageMap map[string]interface{}, format string, message ...interface{}) {
	fields := buildFields(ctx, messageMap)
	staticLogger.WithFields(fields).Infof(format, message...)
}

func Error(ctx context.Context, messageMap map[string]interface{}, message string, err error) {
	fields := buildFields(ctx, messageMap)
	if err != nil {
		fields["error_message"] = err.Error()
	}
	staticLogger.WithFields(fields).Error(message)
}

func Warn(ctx context.Context, messageMap map[string]interface{}, message string, err error) {
	fields := buildFields(ctx, messageMap)
	if err != nil {
		fields["error_message"] = err.Error()
	}
	staticLogger.WithFields(fields).Warn(message)
}

func buildFields(ctx context.Context, messageMap map[string]interface{}) map[string]interface{} {
	fields := map[string]interface{}{}
	for k, v := range messageMap {
		fields[k] = v
	}
	fields["trace_id"] = fmt.Sprint(ctx.Value(XTraceID))
	return fields
}
