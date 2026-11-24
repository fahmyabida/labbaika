package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// Key type for context values
type contextKey string

const XTraceID contextKey = "x-trace-id"

// Fields is extra structured context a caller wants attached to a log line
// (e.g. which record id was involved). Not for identifying which function
// logged — that's auto-captured from the call site.
type Fields = map[string]interface{}

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

func Info(ctx context.Context, message string, extra ...Fields) {
	fields := buildFields(ctx, mergeFields(extra...))
	addCaller(fields)
	staticLogger.WithFields(fields).Info(message)
}

func Infof(ctx context.Context, format string, message ...interface{}) {
	fields := buildFields(ctx, nil)
	addCaller(fields)
	staticLogger.WithFields(fields).Infof(format, message...)
}

// Error logs message/err with the call site (file, line, function) attached
// automatically, so the caller never has to label itself with e.g.
// {"usecase": "cashbook", "method": "Create"}. Pass extra only for business
// context that can't be derived from the stack, like a record id.
func Error(ctx context.Context, message string, err error, extra ...Fields) {
	fields := buildFields(ctx, mergeFields(extra...))
	addCaller(fields)
	if err != nil {
		fields["error_message"] = err.Error()
	}
	staticLogger.WithFields(fields).Error(message)
}

func Warn(ctx context.Context, message string, err error, extra ...Fields) {
	fields := buildFields(ctx, mergeFields(extra...))
	addCaller(fields)
	if err != nil {
		fields["error_message"] = err.Error()
	}
	staticLogger.WithFields(fields).Warn(message)
}

func buildFields(ctx context.Context, extra Fields) Fields {
	fields := Fields{}
	for k, v := range extra {
		fields[k] = v
	}
	fields["trace_id"] = fmt.Sprint(ctx.Value(XTraceID))
	return fields
}

func mergeFields(extra ...Fields) Fields {
	if len(extra) == 0 {
		return nil
	}
	merged := Fields{}
	for _, f := range extra {
		for k, v := range f {
			merged[k] = v
		}
	}
	return merged
}

// addCaller attaches file, line, and function name of whoever called the
// public logger function two frames up (Info/Error/Warn/etc itself is frame 1,
// addCaller is frame 0) — i.e. the real call site in usecase/handler/etc code.
func addCaller(fields Fields) {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return
	}
	fields["file"] = filepath.Base(file)
	fields["line"] = line
	if fn := runtime.FuncForPC(pc); fn != nil {
		name := fn.Name()
		if idx := strings.LastIndex(name, "/"); idx != -1 {
			name = name[idx+1:]
		}
		fields["func"] = name
	}
}
