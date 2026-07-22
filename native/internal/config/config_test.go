package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDevelopmentConfigStartsWithoutRequiredEnvironmentOverrides(t *testing.T) {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve config test path")
	}
	appDir := filepath.Join(filepath.Dir(sourceFile), "..", "..", "cmd", "api")
	t.Setenv(AppEnvVar, "dev")
	t.Setenv("QUICK_ADMIN_DATABASE_DSN", "")
	t.Setenv("QUICK_ADMIN_REDIS_ADDR", "")
	t.Setenv("QUICK_ADMIN_REDIS_PASSWORD", "")
	t.Setenv("QUICK_ADMIN_JWT_SECRET", "")

	cfg, _, err := Load(appDir)
	if err != nil {
		t.Fatalf("Load(dev) error = %v", err)
	}
	if cfg.Database.DSN == "" || cfg.Redis.Addr == "" || cfg.JWT.Secret == "" {
		t.Fatal("development config is missing a required local value")
	}
}

func TestLoadEnvironmentOverridesYAML(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "conf.dev.yaml")
	content := []byte(`
server:
  port: 9009
database:
  dsn: yaml-dsn
redis:
  addr: yaml-redis:6379
jwt:
  secret: yaml-secret-that-is-long-enough-for-tests
cors:
  enabled: false
`)
	if err := os.WriteFile(configPath, content, 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv(AppEnvVar, "dev")
	t.Setenv("QUICK_ADMIN_SERVER_PORT", "8080")
	t.Setenv("QUICK_ADMIN_DATABASE_DSN", "environment-dsn")
	t.Setenv("QUICK_ADMIN_REDIS_ADDR", "environment-redis:6379")
	t.Setenv("QUICK_ADMIN_JWT_SECRET", "environment-secret-with-at-least-32-characters")

	cfg, _, err := Load(dir)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Server.Port != 8080 {
		t.Fatalf("Server.Port = %d, want 8080", cfg.Server.Port)
	}
	if cfg.Database.DSN != "environment-dsn" {
		t.Fatalf("Database.DSN = %q", cfg.Database.DSN)
	}
	if cfg.Redis.Addr != "environment-redis:6379" {
		t.Fatalf("Redis.Addr = %q", cfg.Redis.Addr)
	}
	if cfg.JWT.Secret != "environment-secret-with-at-least-32-characters" {
		t.Fatalf("JWT.Secret was not overridden by the environment")
	}
}

func TestConfigValidate(t *testing.T) {
	t.Parallel()

	valid := Config{
		Server:   Server{Port: 9009},
		Database: Database{DSN: "postgres-dsn"},
		Redis:    Redis{Addr: "redis:6379"},
		JWT:      JWT{Secret: "a-secret-with-at-least-32-characters"},
	}
	if err := valid.Validate("dev"); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}

	tweak := valid
	tweak.CORS.Enabled = true
	if err := tweak.Validate("prod"); err == nil {
		t.Fatal("prod config with CORS enabled was accepted")
	}

	tweak = valid
	tweak.JWT.Secret = "short"
	if err := tweak.Validate("dev"); err == nil {
		t.Fatal("short JWT secret was accepted")
	}
}

func TestCurrentEnvDefaultsToDev(t *testing.T) {
	t.Setenv(AppEnvVar, "")
	if got := CurrentEnv(); got != "dev" {
		t.Fatalf("CurrentEnv() = %q, want dev", got)
	}
}

func TestLoadRejectsUnknownProfile(t *testing.T) {
	t.Setenv(AppEnvVar, "staging")
	if _, _, err := Load(t.TempDir()); err == nil {
		t.Fatal("Load() accepted unsupported profile")
	}
}

func TestEnvironmentName(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"database.maxOpenConns": "QUICK_ADMIN_DATABASE_MAX_OPEN_CONNS",
		"s3.useSSL":             "QUICK_ADMIN_S3_USE_SSL",
		"wechat.appId":          "QUICK_ADMIN_WECHAT_APP_ID",
	}
	for key, want := range tests {
		if got := environmentName(key); got != want {
			t.Errorf("environmentName(%q) = %q, want %q", key, got, want)
		}
	}
}
