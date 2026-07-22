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
	examplePath := filepath.Join(filepath.Dir(sourceFile), "..", "..", "application", "api", "conf.example.yaml")
	example, err := os.ReadFile(examplePath)
	if err != nil {
		t.Fatal(err)
	}
	appDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(appDir, "conf.dev.yaml"), example, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(AppEnvVar, "dev")
	t.Setenv("LIGHTNING_DATABASE_DSN", "")
	t.Setenv("LIGHTNING_REDIS_ADDR", "")
	t.Setenv("LIGHTNING_REDIS_PASSWORD", "")
	t.Setenv("LIGHTNING_JWT_SECRET", "")

	cfg, _, err := Load(appDir, ServiceAPI)
	if err != nil {
		t.Fatalf("Load(dev) error = %v", err)
	}
	if cfg.Database.DSN == "" || cfg.Redis.Addr == "" || cfg.JWT.Secret == "" {
		t.Fatal("development config is missing a required local value")
	}
	if cfg.AppDir != appDir {
		t.Fatalf("AppDir = %q, want %q", cfg.AppDir, appDir)
	}
}

func TestServiceConfigExamplesSupportBothProfiles(t *testing.T) {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve config test path")
	}
	nativeDir := filepath.Join(filepath.Dir(sourceFile), "..", "..")
	services := []struct {
		name    string
		service Service
	}{
		{name: "api", service: ServiceAPI},
		{name: "scheduler", service: ServiceScheduler},
	}
	for _, tt := range services {
		t.Run(tt.name, func(t *testing.T) {
			example, err := os.ReadFile(filepath.Join(nativeDir, "application", tt.name, "conf.example.yaml"))
			if err != nil {
				t.Fatal(err)
			}
			for _, profile := range []string{"dev", "prod"} {
				profile := profile
				t.Run(profile, func(t *testing.T) {
					dir := t.TempDir()
					if err := os.WriteFile(filepath.Join(dir, "conf."+profile+".yaml"), example, 0o600); err != nil {
						t.Fatal(err)
					}
					t.Setenv(AppEnvVar, profile)
					cfg, _, err := Load(dir, tt.service)
					if err != nil {
						t.Fatalf("Load(%s/%s) error = %v", tt.name, profile, err)
					}
					if cfg.Database.DSN == "" {
						t.Fatalf("%s example is missing database.dsn", tt.name)
					}
					if tt.service == ServiceAPI && (cfg.Redis.Addr == "" || cfg.RabbitMQ.URL == "" || cfg.S3.Endpoint == "") {
						t.Fatal("API example contains an empty infrastructure default")
					}
				})
			}
		})
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
	t.Setenv("LIGHTNING_SERVER_PORT", "8080")
	t.Setenv("LIGHTNING_DATABASE_DSN", "environment-dsn")
	t.Setenv("LIGHTNING_REDIS_ADDR", "environment-redis:6379")
	t.Setenv("LIGHTNING_JWT_SECRET", "environment-secret-with-at-least-32-characters")

	cfg, _, err := Load(dir, ServiceAPI)
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
	if err := valid.Validate("dev", ServiceAPI); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}

	tweak := valid
	tweak.CORS.Enabled = true
	if err := tweak.Validate("prod", ServiceAPI); err == nil {
		t.Fatal("prod config with CORS enabled was accepted")
	}

	tweak = valid
	tweak.JWT.Secret = "short"
	if err := tweak.Validate("dev", ServiceAPI); err == nil {
		t.Fatal("short JWT secret was accepted")
	}

	scheduler := Config{Database: Database{DSN: "postgres-dsn"}}
	if err := scheduler.Validate("prod", ServiceScheduler); err != nil {
		t.Fatalf("minimal scheduler config rejected: %v", err)
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
	if _, _, err := Load(t.TempDir(), ServiceAPI); err == nil {
		t.Fatal("Load() accepted unsupported profile")
	}
}

func TestEnvironmentName(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"database.maxOpenConns": "LIGHTNING_DATABASE_MAX_OPEN_CONNS",
		"s3.useSSL":             "LIGHTNING_S3_USE_SSL",
	}
	for key, want := range tests {
		if got := environmentName(key); got != want {
			t.Errorf("environmentName(%q) = %q, want %q", key, got, want)
		}
	}
}
