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
	addTraceID(ctx, messageMap)
	staticLogger.WithFields(messageMap).Info(message)
}

func Infof(ctx context.Context, messageMap map[string]interface{}, format string, message ...interface{}) {
	addTraceID(ctx, messageMap)
	staticLogger.WithFields(messageMap).Infof(format, message...)
}

func Error(ctx context.Context, messageMap map[string]interface{}, message string, err error) {
	addTraceID(ctx, messageMap)
	messageMap["error_message"] = err.Error()
	staticLogger.WithFields(messageMap).Error(message)
}

func addTraceID(ctx context.Context, messageMap map[string]interface{}) {
	traceID := ctx.Value(XTraceID)
	initMessageMap(&messageMap)
	messageMap["trace_id"] = fmt.Sprint(traceID)
}

func initMessageMap(messageMap *map[string]interface{}) {
	if messageMap == nil {
		messageMap = &map[string]interface{}{}
	}
}
