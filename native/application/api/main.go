package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gcc798/lightning/application/api/bootstrap"
	"github.com/gcc798/lightning/application/api/httpx"
	"github.com/gcc798/lightning/application/api/middleware"
	"github.com/gcc798/lightning/application/api/migrations"
	_ "github.com/gcc798/lightning/application/api/openapi" // 注册 Swagger 文档
	"github.com/gcc798/lightning/application/api/router"
	"github.com/gcc798/lightning/application/api/validator"
	"github.com/gcc798/lightning/internal/app"
	"github.com/gcc798/lightning/internal/config"
	"github.com/gcc798/lightning/internal/container"
	"github.com/gcc798/lightning/internal/modules"
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"go.uber.org/zap"
)

//	@title			NTZ API 文档
//	@version		1.0
//	@description	Quick Admin RESTful API 接口文档，支持双 Token 认证机制
//	@termsOfService	https://example.com/terms/

//	@contact.name	技术支持
//	@contact.url	https://example.com/support
//	@contact.email	support@example.com

//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT

//	@host		localhost:9009
//	@BasePath	/

//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				格式: "Bearer {access_token}"，AccessToken 从登录接口获取

//	@tag.name			认证
//	@tag.description	用户认证相关接口，包括登录、登出、Token 刷新等

//	@tag.name			用户管理
//	@tag.description	用户信息管理接口

func main() {
	base, err := app.New("application/api", config.ServiceAPI, container.WithAPIInfrastructure())
	if err != nil {
		fmt.Printf("failed to initialize API process: %v\n", err)
		os.Exit(1)
	}
	c, cfg := base.Container, base.Config

	logger := c.GetLogger()
	sqlDB, err := c.GetDB().DB()
	if err != nil {
		logger.Fatal("failed to access database connection", zap.Error(err))
	}
	logger.Info("starting API database migrations")
	if err := migrations.Up(sqlDB); err != nil {
		logger.Fatal("failed to migrate API database", zap.Error(err))
	}
	logger.Info("API database migrations completed")
	logger.Info("container initialized successfully")
	if err := c.RegisterModules(context.Background(),
		modules.NewSMSModule(),
		modules.NewEmailModule(),
		modules.NewWeChatModule(),
		modules.NewCaptchaModule(),
	); err != nil {
		logger.Fatal("failed to initialize API modules", zap.Error(err))
	}

	b, err := bootstrap.New(c)
	if err != nil {
		logger.Fatal("failed to bootstrap business components", zap.Error(err))
	}

	if err := c.Start(); err != nil {
		logger.Fatal("failed to start API infrastructure", zap.Error(err))
	}
	if err := c.StartModules(context.Background()); err != nil {
		logger.Fatal("failed to start API modules", zap.Error(err))
	}

	// 初始化HTTP引擎
	e := echo.New()
	//binding.EnableDecoderUseNumber = true

	// 初始化中文验证错误翻译器
	validator.Init()
	e.Validator = validator.EchoValidator{}
	e.Binder = &validator.EchoBinder{}

	// 使用自定义的 Recovery 中间件（支持结构化错误处理）
	// 自定义 Recovery 提供结构化日志和统一错误响应。
	e.Use(middleware.Recovery(logger.Get()))
	e.Use(echoMiddleware.RequestLogger())

	// 本地前后端分离开发时启用；生产环境通过同源 Nginx 代理，不启用应用层 CORS。
	if cfg.CORS.Enabled {
		e.Use(middleware.CORS())
	}

	// 添加字符串ID转换中间件（处理前端传递的字符串ID）
	e.Use(middleware.StringIDConverter())

	// 添加操作日志中间件（全局记录所有接口访问）
	operLogWriter := middleware.NewOperLogWriter(c.GetDB(), c.GetLogger())
	e.Use(middleware.OperationLog(operLogWriter))

	// 注册路由
	if err := router.Setup(httpx.NewRouter(e), c, b); err != nil {
		logger.Fatal("failed to register routes", zap.Error(err))
	}

	// 配置HTTP服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:           addr,
		Handler:        e,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// 在 goroutine 中启动服务器
	go func() {
		logger.Info("starting http server", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("failed to start http server", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	// 默认终止命令发送 syscall.SIGTERM。
	// 中断命令发送 syscall.SIGINT。
	// 强制终止命令发送 syscall.SIGKILL，进程无法捕获，不需要监听。
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server...")

	// 设置 30 秒的超时时间用于优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 优雅关闭 HTTP 服务器
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", zap.Error(err))
	}
	operLogWriter.Stop()
	if err := c.StopModules(ctx); err != nil {
		logger.Error("failed to stop API modules", zap.Error(err))
	}
	c.Stop()

	logger.Info("server exited gracefully")
}
