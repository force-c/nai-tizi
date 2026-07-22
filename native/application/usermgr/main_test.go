package main

import "testing"

func TestParseOptionsAndValidateInput(t *testing.T) {
	opts, err := parseOptions([]string{
		"--operation=create",
		"--username=admin",
		"--nickname=管理员",
		"--role=super_admin",
		"--config-dir=/app",
	})
	if err != nil {
		t.Fatalf("parseOptions() error = %v", err)
	}
	if opts.operation != "create" || opts.username != "admin" || opts.configDir != "/app" {
		t.Fatalf("parseOptions() = %#v", opts)
	}
	if err := validateInput(opts, "strong-password"); err != nil {
		t.Fatalf("validateInput() error = %v", err)
	}
}

func TestValidateInputRejectsUnsafeInput(t *testing.T) {
	tests := []struct {
		name     string
		opts     options
		password string
	}{
		{name: "unknown operation", opts: options{operation: "delete", username: "admin"}, password: "strong-password"},
		{name: "missing username", opts: options{operation: "reset"}, password: "strong-password"},
		{name: "short password", opts: options{operation: "reset", username: "admin"}, password: "short"},
		{name: "missing role", opts: options{operation: "create", username: "admin", nickname: "管理员"}, password: "strong-password"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateInput(tt.opts, tt.password); err == nil {
				t.Fatal("validateInput() expected error")
			}
		})
	}
}
