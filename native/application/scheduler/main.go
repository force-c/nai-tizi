package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gcc798/lightning/application/scheduler/jobs"
	"github.com/gcc798/lightning/internal/app"
	"github.com/gcc798/lightning/internal/config"
	"github.com/gcc798/lightning/internal/container"
	"github.com/gcc798/lightning/internal/modules"
	"go.uber.org/zap"
)

func main() {
	base, err := app.New("application/scheduler", config.ServiceScheduler)
	if err != nil {
		fmt.Printf("failed to initialize scheduler process: %v\n", err)
		os.Exit(1)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, base.Container); err != nil {
		base.Logger.Error("scheduler process exited with error", zap.Error(err))
		os.Exit(1)
	}
}

func run(ctx context.Context, cont container.Container) error {
	definitions := jobs.Definitions(cont.GetDB(), cont.GetLogger())
	if err := cont.RegisterModules(ctx,
		modules.NewSchedulerModule(definitions),
	); err != nil {
		return fmt.Errorf("initialize scheduler modules: %w", err)
	}
	if err := cont.Start(); err != nil {
		return fmt.Errorf("start scheduler infrastructure: %w", err)
	}
	if err := cont.StartModules(ctx); err != nil {
		cont.Stop()
		return fmt.Errorf("start scheduler modules: %w", err)
	}
	cont.GetLogger().Info("scheduler process started")
	<-ctx.Done()
	cont.GetLogger().Info("stopping scheduler process")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	moduleErr := cont.StopModules(shutdownCtx)
	cont.Stop()
	return errors.Join(moduleErr)
}
