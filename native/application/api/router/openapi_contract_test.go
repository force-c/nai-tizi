package router

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestOpenAPICriticalOperations(t *testing.T) {
	t.Parallel()

	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve test source path")
	}
	documentPath := filepath.Join(filepath.Dir(sourceFile), "..", "openapi", "swagger.json")
	data, err := os.ReadFile(documentPath)
	if err != nil {
		t.Fatal(err)
	}
	var document struct {
		Paths map[string]map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatal(err)
	}

	operations := map[string][]string{
		"/login":        {"post"},
		"/logout":       {"post"},
		"/auth/refresh": {"post"},
		"/health/live":  {"get"},
		"/health/ready": {"get"},
		"/api/v1/user":  {"post"},
		"/api/v1/role":  {"post"},
		"/api/v1/menu":  {"get", "post"},
	}
	for path, methods := range operations {
		pathItem, exists := document.Paths[path]
		if !exists {
			t.Errorf("OpenAPI path %s is missing", path)
			continue
		}
		for _, method := range methods {
			if _, exists := pathItem[method]; !exists {
				t.Errorf("OpenAPI operation %s %s is missing", method, path)
			}
		}
	}
}
