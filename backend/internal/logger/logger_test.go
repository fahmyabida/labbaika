package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestError_AttachesCallSite(t *testing.T) {
	InitStaticLogger("test")
	var buf bytes.Buffer
	staticLogger.SetOutput(&buf)

	Error(context.Background(), "boom happened", errors.New("db down"))

	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("invalid log json: %v", err)
	}

	if entry["file"] != "logger_test.go" {
		t.Errorf("file = %v, want logger_test.go (not logger.go — would mean the caller skip-depth is wrong)", entry["file"])
	}
	if fn, _ := entry["func"].(string); !strings.Contains(fn, "TestError_AttachesCallSite") {
		t.Errorf("func = %v, want to contain TestError_AttachesCallSite", fn)
	}
	if entry["error_message"] != "db down" {
		t.Errorf("error_message = %v, want %q", entry["error_message"], "db down")
	}
}
