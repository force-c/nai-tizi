// Package app provides reusable process bootstrap for application entrypoints.
package app

import (
	"fmt"

	"github.com/gcc798/lightning/internal/config"
	"github.com/gcc798/lightning/internal/container"
	logging "github.com/gcc798/lightning/internal/logger"
	"github.com/spf13/viper"
)

// Base contains the shared infrastructure of one native process.
type Base struct {
	Config    *config.Config
	Viper     *viper.Viper
	Logger    logging.Logger
	Container container.Container
}

// New initializes configuration, logging, database and Redis, then applies
// process-specific container options.
func New(configDir string, service config.Service, options ...container.Option) (*Base, error) {
	cfg, v, err := config.Load(configDir, service)
	if err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}
	log, err := logging.NewLogger(config.CurrentEnv(), cfg.AppDir)
	if err != nil {
		return nil, fmt.Errorf("initialize logger: %w", err)
	}
	cont, err := container.New(cfg, v, log, options...)
	if err != nil {
		return nil, fmt.Errorf("initialize container: %w", err)
	}
	return &Base{Config: cfg, Viper: v, Logger: log, Container: cont}, nil
}
