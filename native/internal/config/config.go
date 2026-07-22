package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/spf13/viper"
)

const AppEnvVar = "LIGHTNING_APP_ENV"

type Service string

const (
	ServiceAPI         Service = "api"
	ServiceScheduler   Service = "scheduler"
	ServiceUserManager Service = "usermgr"
)

type Server struct {
	Port int `mapstructure:"port"`
}

type Database struct {
	DSN                    string `mapstructure:"dsn"`
	MaxOpenConns           int    `mapstructure:"maxOpenConns"`
	MaxIdleConns           int    `mapstructure:"maxIdleConns"`
	ConnMaxLifetimeMinutes int    `mapstructure:"connMaxLifetimeMinutes"`
	SlowThreshold          int    `mapstructure:"slowThreshold"`
}

type Redis struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWT struct {
	Secret string `mapstructure:"secret"`
	Expire int64  `mapstructure:"expire"`
}

type Auth struct {
	TokenHeader     string `mapstructure:"tokenHeader"`
	AllowConcurrent bool   `mapstructure:"allowConcurrent"`
}

type CORS struct {
	Enabled bool `mapstructure:"enabled"`
}

type RabbitMQ struct {
	Enabled bool   `mapstructure:"enabled"`
	URL     string `mapstructure:"url"`
}

type S3 struct {
	Enabled         bool   `mapstructure:"enabled"`
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"accessKeyId"`
	SecretAccessKey string `mapstructure:"secretAccessKey"`
	Region          string `mapstructure:"region"`
	Bucket          string `mapstructure:"bucket"`
	UseSSL          bool   `mapstructure:"useSSL"`
	ForcePathStyle  bool   `mapstructure:"forcePathStyle"`
}

type WebSocket struct {
	Enabled             bool `mapstructure:"enabled"`
	TimeoutEnabled      bool `mapstructure:"timeoutEnabled"`
	ReadTimeoutSeconds  int  `mapstructure:"readTimeoutSeconds"`
	WriteTimeoutSeconds int  `mapstructure:"writeTimeoutSeconds"`
	HeartbeatEnabled    bool `mapstructure:"heartbeatEnabled"`
	MaxReadTimeouts     int  `mapstructure:"maxReadTimeouts"`
}

type Config struct {
	AppDir    string    `mapstructure:"-"`
	Server    Server    `mapstructure:"server"`
	Database  Database  `mapstructure:"database"`
	Redis     Redis     `mapstructure:"redis"`
	JWT       JWT       `mapstructure:"jwt"`
	Auth      Auth      `mapstructure:"auth"`
	CORS      CORS      `mapstructure:"cors"`
	RabbitMQ  RabbitMQ  `mapstructure:"rabbitmq"`
	S3        S3        `mapstructure:"s3"`
	WebSocket WebSocket `mapstructure:"websocket"`
}

func Load(configDir string, service Service) (*Config, *viper.Viper, error) {
	profile := CurrentEnv()
	if profile != "dev" && profile != "prod" {
		return nil, nil, fmt.Errorf("%s must be dev or prod", AppEnvVar)
	}
	if service != ServiceAPI && service != ServiceScheduler && service != ServiceUserManager {
		return nil, nil, fmt.Errorf("unknown config service %q", service)
	}
	v := viper.New()
	setDefaults(v)
	v.SetEnvPrefix("LIGHTNING")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	if err := bindEnvironment(v); err != nil {
		return nil, nil, err
	}

	configFileName := fmt.Sprintf("conf.%s.yaml", profile)
	foundPath, err := ResolveFilePath(configDir, configFileName)
	if err != nil {
		return nil, nil, err
	}
	v.SetConfigFile(foundPath)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return nil, nil, fmt.Errorf("read config from %s: %w", foundPath, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, nil, fmt.Errorf("decode config: %w", err)
	}
	cfg.AppDir = filepath.Dir(foundPath)
	if err := cfg.Validate(profile, service); err != nil {
		return nil, nil, err
	}
	return &cfg, v, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.port", 9009)
	v.SetDefault("database.maxOpenConns", 100)
	v.SetDefault("database.maxIdleConns", 10)
	v.SetDefault("database.connMaxLifetimeMinutes", 60)
	v.SetDefault("database.slowThreshold", 200)
	v.SetDefault("redis.db", 0)
	v.SetDefault("jwt.expire", 7200)
	v.SetDefault("auth.tokenHeader", "Authorization")
	v.SetDefault("auth.allowConcurrent", false)
	v.SetDefault("cors.enabled", false)
}

func bindEnvironment(v *viper.Viper) error {
	keys := []string{
		"server.port",
		"database.dsn", "database.maxOpenConns", "database.maxIdleConns", "database.connMaxLifetimeMinutes", "database.slowThreshold",
		"redis.addr", "redis.password", "redis.db",
		"jwt.secret", "jwt.expire",
		"auth.tokenHeader", "auth.allowConcurrent",
		"cors.enabled",
		"rabbitmq.enabled", "rabbitmq.url",
		"s3.enabled", "s3.endpoint", "s3.accessKeyId", "s3.secretAccessKey", "s3.region", "s3.bucket", "s3.useSSL", "s3.forcePathStyle",
		"websocket.enabled", "websocket.timeoutEnabled", "websocket.readTimeoutSeconds", "websocket.writeTimeoutSeconds", "websocket.heartbeatEnabled", "websocket.maxReadTimeouts",
	}
	for _, key := range keys {
		if err := v.BindEnv(key, environmentName(key)); err != nil {
			return fmt.Errorf("bind environment for %s: %w", key, err)
		}
	}
	return nil
}

func environmentName(key string) string {
	runes := []rune(key)
	var name strings.Builder
	name.WriteString("LIGHTNING_")
	for i, current := range runes {
		switch {
		case current == '.':
			name.WriteByte('_')
		case unicode.IsUpper(current):
			previousIsLower := i > 0 && (unicode.IsLower(runes[i-1]) || unicode.IsDigit(runes[i-1]))
			startsWord := i > 0 && unicode.IsUpper(runes[i-1]) && i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if previousIsLower || startsWord {
				name.WriteByte('_')
			}
			name.WriteRune(current)
		default:
			name.WriteRune(unicode.ToUpper(current))
		}
	}
	return name.String()
}

func (c *Config) Validate(profile string, service Service) error {
	if c.Database.DSN == "" {
		return fmt.Errorf("database.dsn is required")
	}
	if service != ServiceAPI {
		return nil
	}
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("server.port must be between 1 and 65535")
	}
	if c.Redis.Addr == "" {
		return fmt.Errorf("redis.addr is required")
	}
	if len(c.JWT.Secret) < 32 {
		return fmt.Errorf("jwt.secret must contain at least 32 characters")
	}
	if profile == "prod" && c.CORS.Enabled {
		return fmt.Errorf("cors must be disabled in prod")
	}
	if c.RabbitMQ.Enabled && c.RabbitMQ.URL == "" {
		return fmt.Errorf("rabbitmq.url is required when RabbitMQ is enabled")
	}
	if c.S3.Enabled && (c.S3.Endpoint == "" || c.S3.AccessKeyID == "" || c.S3.SecretAccessKey == "" || c.S3.Bucket == "") {
		return fmt.Errorf("S3 endpoint, credentials and bucket are required when S3 is enabled")
	}
	return nil
}

func CurrentEnv() string {
	if profile := os.Getenv(AppEnvVar); profile != "" {
		return profile
	}
	return "dev"
}

func ResolveFilePath(configDir, fileName string) (string, error) {
	if strings.TrimSpace(configDir) == "" {
		return "", errors.New("config directory is required")
	}
	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("resolve working directory: %w", err)
	}
	executable, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("resolve executable: %w", err)
	}

	var candidates []string
	if filepath.IsAbs(configDir) {
		candidates = []string{filepath.Join(configDir, fileName)}
	} else {
		candidates = []string{
			filepath.Join(workDir, configDir, fileName),
			filepath.Join(filepath.Dir(executable), fileName),
		}
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("config file not found: tried %v", candidates)
}
