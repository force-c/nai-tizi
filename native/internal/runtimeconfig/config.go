// Package runtimeconfig owns database-backed module configuration contracts.
package runtimeconfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/robfig/cron/v3"
)

const (
	CodeWeChat    = "integration.wechat"
	CodeSMS       = "integration.sms"
	CodeEmail     = "integration.email"
	CodeCaptcha   = "auth.captcha"
	CodeScheduler = "scheduler"
)

// WeChatConfig configures the WeChat capability.
type WeChatConfig struct {
	Enabled    bool   `json:"enabled"`
	AppID      string `json:"appId"`
	Secret     string `json:"secret"`
	TemplateID string `json:"templateId"`
}

// SMSConfig configures the SMS capability.
type SMSConfig struct {
	Enabled         bool   `json:"enabled"`
	AccessKeyID     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	SignName        string `json:"signName"`
	TemplateCode    string `json:"templateCode"`
}

// EmailConfig configures the email capability.
type EmailConfig struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
}

// CaptchaConfig configures all captcha providers.
type CaptchaConfig struct {
	Image ImageCaptchaConfig `json:"image"`
	SMS   SMSCaptchaConfig   `json:"sms"`
	Email EmailCaptchaConfig `json:"email"`
}

type ImageCaptchaConfig struct {
	Enabled bool `json:"enabled"`
	Length  int  `json:"length"`
	Width   int  `json:"width"`
	Height  int  `json:"height"`
	Expire  int  `json:"expire"`
}

type SMSCaptchaConfig struct {
	Enabled  bool   `json:"enabled"`
	Length   int    `json:"length"`
	Expire   int    `json:"expire"`
	Template string `json:"template"`
	Provider string `json:"provider"`
}

type EmailCaptchaConfig struct {
	Enabled  bool   `json:"enabled"`
	Length   int    `json:"length"`
	Expire   int    `json:"expire"`
	Template string `json:"template"`
}

// SchedulerConfig configures code-registered jobs.
type SchedulerConfig struct {
	Enabled                bool                 `json:"enabled"`
	RefreshIntervalSeconds int                  `json:"refreshIntervalSeconds"`
	Jobs                   map[string]JobConfig `json:"jobs"`
}

type JobConfig struct {
	Enabled bool   `json:"enabled"`
	Cron    string `json:"cron"`
}

// Validate validates JSON for known runtime module configuration codes.
func Validate(code string, data []byte) error {
	if len(data) == 0 || !json.Valid(data) {
		return errors.New("configuration data must be valid JSON")
	}
	switch code {
	case CodeWeChat:
		var cfg WeChatConfig
		if err := decodeStrict(data, &cfg); err != nil {
			return err
		}
		if cfg.Enabled && (strings.TrimSpace(cfg.AppID) == "" || strings.TrimSpace(cfg.Secret) == "") {
			return errors.New("wechat appId and secret are required when enabled")
		}
	case CodeSMS:
		var cfg SMSConfig
		if err := decodeStrict(data, &cfg); err != nil {
			return err
		}
		if cfg.Enabled && (cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" || cfg.SignName == "" || cfg.TemplateCode == "") {
			return errors.New("SMS credentials, signName and templateCode are required when enabled")
		}
	case CodeEmail:
		var cfg EmailConfig
		if err := decodeStrict(data, &cfg); err != nil {
			return err
		}
		if cfg.Enabled && (cfg.Host == "" || cfg.Port <= 0 || cfg.Username == "" || cfg.Password == "" || cfg.From == "") {
			return errors.New("email host, port, username, password and from are required when enabled")
		}
	case CodeCaptcha:
		var cfg CaptchaConfig
		if err := decodeStrict(data, &cfg); err != nil {
			return err
		}
		if cfg.Image.Enabled && (cfg.Image.Length <= 0 || cfg.Image.Width <= 0 || cfg.Image.Height <= 0 || cfg.Image.Expire <= 0) {
			return errors.New("enabled image captcha requires positive length, width, height and expire")
		}
		if cfg.SMS.Enabled && (cfg.SMS.Length <= 0 || cfg.SMS.Expire <= 0 || cfg.SMS.Template == "" || cfg.SMS.Provider == "") {
			return errors.New("enabled SMS captcha requires positive length and expire, template and provider")
		}
		if cfg.Email.Enabled && (cfg.Email.Length <= 0 || cfg.Email.Expire <= 0 || cfg.Email.Template == "") {
			return errors.New("enabled email captcha requires positive length and expire and template")
		}
	case CodeScheduler:
		var cfg SchedulerConfig
		if err := decodeStrict(data, &cfg); err != nil {
			return err
		}
		if cfg.Enabled && cfg.RefreshIntervalSeconds <= 0 {
			return errors.New("enabled scheduler requires a positive refreshIntervalSeconds")
		}
		parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		for name, job := range cfg.Jobs {
			if strings.TrimSpace(name) == "" {
				return errors.New("scheduler job name cannot be empty")
			}
			if job.Enabled {
				if _, err := parser.Parse(job.Cron); err != nil {
					return fmt.Errorf("scheduler job %q has invalid cron: %w", name, err)
				}
			}
		}
	}
	return nil
}

func decodeStrict(data []byte, target any) error {
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode runtime configuration: %w", err)
	}
	return nil
}
