package logger

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestServiceLoggerExamples(t *testing.T) {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve logger test path")
	}
	nativeDir := filepath.Join(filepath.Dir(sourceFile), "..", "..")
	for _, service := range []string{"api", "scheduler"} {
		cfg, err := LoadConfig(filepath.Join(nativeDir, "application", service, "zaplogger.example.yaml"))
		if err != nil {
			t.Fatalf("LoadConfig(%s) error = %v", service, err)
		}
		if cfg.Level == "" || cfg.Output == "" || cfg.Encoding == "" || cfg.File.Filename == "" {
			t.Fatalf("%s logger example contains an empty required value", service)
		}
	}
}
