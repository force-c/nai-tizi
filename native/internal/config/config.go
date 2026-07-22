package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/spf13/viper"
)

const AppEnvVar = "QUICK_ADMIN_APP_ENV"

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

type WeChat struct {
	Enabled    bool   `mapstructure:"enabled" json:"enabled"`
	AppID      string `mapstructure:"appId" json:"appid"`
	Secret     string `mapstructure:"secret" json:"secret"`
	TemplateID string `mapstructure:"templateId" json:"templateId"`
}

type MQTT struct {
	Enabled  bool   `mapstructure:"enabled"`
	Broker   string `mapstructure:"broker"`
	ClientID string `mapstructure:"clientId"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	QoS      byte   `mapstructure:"qos"`
}

type RabbitMQ struct {
	Enabled bool   `mapstructure:"enabled"`
	URL     string `mapstructure:"url"`
}

type SMS struct {
	Enabled         bool   `mapstructure:"enabled"`
	AccessKeyId     string `mapstructure:"accessKeyId"`
	AccessKeySecret string `mapstructure:"accessKeySecret"`
	SignName        string `mapstructure:"signName"`
	TemplateCode    string `mapstructure:"templateCode"`
}

type Email struct {
	Enabled  bool   `mapstructure:"enabled"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
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

type Scheduler struct {
	Enabled bool `mapstructure:"enabled"`
}

type WebSocket struct {
	Enabled             bool `mapstructure:"enabled"`
	TimeoutEnabled      bool `mapstructure:"timeoutEnabled"`
	ReadTimeoutSeconds  int  `mapstructure:"readTimeoutSeconds"`
	WriteTimeoutSeconds int  `mapstructure:"writeTimeoutSeconds"`
	HeartbeatEnabled    bool `mapstructure:"heartbeatEnabled"`
	MaxReadTimeouts     int  `mapstructure:"maxReadTimeouts"`
}

type Captcha struct {
	Image ImageCaptcha `mapstructure:"image"`
	SMS   SMSCaptcha   `mapstructure:"sms"`
	Email EmailCaptcha `mapstructure:"email"`
}

type ImageCaptcha struct {
	Enabled bool `mapstructure:"enabled"`
	Length  int  `mapstructure:"length"`
	Width   int  `mapstructure:"width"`
	Height  int  `mapstructure:"height"`
	Expire  int  `mapstructure:"expire"`
}

type SMSCaptcha struct {
	Enabled  bool   `mapstructure:"enabled"`
	Length   int    `mapstructure:"length"`
	Expire   int    `mapstructure:"expire"`
	Template string `mapstructure:"template"`
	Provider string `mapstructure:"provider"`
}

type EmailCaptcha struct {
	Enabled  bool   `mapstructure:"enabled"`
	Length   int    `mapstructure:"length"`
	Expire   int    `mapstructure:"expire"`
	Template string `mapstructure:"template"`
}

type Config struct {
	AppDir    string    `mapstructure:"-"`
	Server    Server    `mapstructure:"server"`
	Database  Database  `mapstructure:"database"`
	Redis     Redis     `mapstructure:"redis"`
	JWT       JWT       `mapstructure:"jwt"`
	Auth      Auth      `mapstructure:"auth"`
	CORS      CORS      `mapstructure:"cors"`
	Captcha   Captcha   `mapstructure:"captcha"`
	WeChat    WeChat    `mapstructure:"wechat"`
	MQTT      MQTT      `mapstructure:"mqtt"`
	RabbitMQ  RabbitMQ  `mapstructure:"rabbitmq"`
	SMS       SMS       `mapstructure:"sms"`
	Email     Email     `mapstructure:"email"`
	S3        S3        `mapstructure:"s3"`
	Scheduler Scheduler `mapstructure:"scheduler"`
	WebSocket WebSocket `mapstructure:"websocket"`
}

func Load(appDir string) (*Config, *viper.Viper, error) {
	profile := CurrentEnv()
	if profile != "dev" && profile != "prod" {
		return nil, nil, fmt.Errorf("%s must be dev or prod", AppEnvVar)
	}
	v := viper.New()
	setDefaults(v)
	v.SetEnvPrefix("QUICK_ADMIN")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	if err := bindEnvironment(v); err != nil {
		return nil, nil, err
	}

	configFileName := fmt.Sprintf("conf.%s.yaml", profile)
	foundPath, err := ResolveFilePath(appDir, configFileName)
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
	cfg.AppDir = appDir
	if err := cfg.Validate(profile); err != nil {
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
	v.SetDefault("captcha.image.length", 4)
	v.SetDefault("captcha.image.width", 120)
	v.SetDefault("captcha.image.height", 40)
	v.SetDefault("captcha.image.expire", 300)
	v.SetDefault("captcha.sms.length", 6)
	v.SetDefault("captcha.sms.expire", 300)
	v.SetDefault("captcha.email.length", 6)
	v.SetDefault("captcha.email.expire", 300)
}

func bindEnvironment(v *viper.Viper) error {
	keys := []string{
		"server.port",
		"database.dsn", "database.maxOpenConns", "database.maxIdleConns", "database.connMaxLifetimeMinutes", "database.slowThreshold",
		"redis.addr", "redis.password", "redis.db",
		"jwt.secret", "jwt.expire",
		"auth.tokenHeader", "auth.allowConcurrent",
		"cors.enabled",
		"captcha.image.enabled", "captcha.image.length", "captcha.image.width", "captcha.image.height", "captcha.image.expire",
		"captcha.sms.enabled", "captcha.sms.length", "captcha.sms.expire", "captcha.sms.template", "captcha.sms.provider",
		"captcha.email.enabled", "captcha.email.length", "captcha.email.expire", "captcha.email.template",
		"wechat.enabled", "wechat.appId", "wechat.secret", "wechat.templateId",
		"mqtt.enabled", "mqtt.broker", "mqtt.clientId", "mqtt.username", "mqtt.password", "mqtt.qos",
		"rabbitmq.enabled", "rabbitmq.url",
		"sms.enabled", "sms.accessKeyId", "sms.accessKeySecret", "sms.signName", "sms.templateCode",
		"email.enabled", "email.host", "email.port", "email.username", "email.password", "email.from",
		"s3.enabled", "s3.endpoint", "s3.accessKeyId", "s3.secretAccessKey", "s3.region", "s3.bucket", "s3.useSSL", "s3.forcePathStyle",
		"scheduler.enabled",
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
	name.WriteString("QUICK_ADMIN_")
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

func (c *Config) Validate(profile string) error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("server.port must be between 1 and 65535")
	}
	if c.Database.DSN == "" {
		return fmt.Errorf("database.dsn is required")
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
	if c.MQTT.Enabled && (c.MQTT.Broker == "" || c.MQTT.ClientID == "") {
		return fmt.Errorf("mqtt.broker and mqtt.clientId are required when MQTT is enabled")
	}
	if c.RabbitMQ.Enabled && c.RabbitMQ.URL == "" {
		return fmt.Errorf("rabbitmq.url is required when RabbitMQ is enabled")
	}
	if c.WeChat.Enabled && (c.WeChat.AppID == "" || c.WeChat.Secret == "") {
		return fmt.Errorf("wechat.appId and wechat.secret are required when WeChat is enabled")
	}
	if c.SMS.Enabled && (c.SMS.AccessKeyId == "" || c.SMS.AccessKeySecret == "" || c.SMS.SignName == "" || c.SMS.TemplateCode == "") {
		return fmt.Errorf("SMS credentials, signName and templateCode are required when SMS is enabled")
	}
	if c.Email.Enabled && (c.Email.Host == "" || c.Email.Port <= 0 || c.Email.Username == "" || c.Email.Password == "" || c.Email.From == "") {
		return fmt.Errorf("email host, port, username, password and from are required when email is enabled")
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

func ResolveFilePath(appDir, fileName string) (string, error) {
	workDir, _ := os.Getwd()
	candidates := []string{
		filepath.Join(workDir, fileName),
		filepath.Join(workDir, "cmd", "api", fileName),
		filepath.Join(appDir, fileName),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("config file not found: tried %v", candidates)
}
