package router

import (
	"github.com/gcc798/lightning/application/api/bootstrap"
	"github.com/gcc798/lightning/application/api/httpx"
	"github.com/gcc798/lightning/application/api/middleware"
	"github.com/gcc798/lightning/application/api/service"
	"github.com/gcc798/lightning/internal/container"
	"github.com/labstack/echo/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// RouterContext 定义业务数据结构。
type RouterContext struct {
	Container         container.Container
	Bootstrap         *bootstrap.Bootstrap
	TokenManager      service.TokenManager
	PermissionService service.PermissionService
	AuthMiddleware    echo.MiddlewareFunc
}

// Setup 配置所有路由。
func Setup(r *httpx.Router, c container.Container, b *bootstrap.Bootstrap) error {
	// 添加 Prometheus 指标收集中间件
	r.Use(middleware.PrometheusMiddleware())

	// 指标端点
	r.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// 接口文档
	r.GET("/swagger/*", echo.WrapHandler(httpSwagger.Handler()))

	// 初始化统一的中间件（除了 auth 模块，其他模块都需要认证）
	tokenManager := service.NewTokenManager(c.GetJWT(), c.GetRedis(), c.GetLogger())
	permissionService := service.NewPermissionService(c.GetDB(), c.GetLogger())
	authMiddleware := middleware.Auth(tokenManager, c.GetConfig())

	// 创建路由上下文
	ctx := &RouterContext{
		Container:         c,
		Bootstrap:         b,
		TokenManager:      tokenManager,
		PermissionService: permissionService,
		AuthMiddleware:    authMiddleware,
	}

	// 注册公共路由（无前缀，部分需要认证）
	registerCommonRoutes(r, ctx)

	// 注册认证相关路由（部分公开，部分需要认证）
	if err := registerAuthRoutes(r, ctx); err != nil {
		return err
	}

	// 注册验证码路由（公开）
	if err := registerCaptchaRoutes(r, ctx); err != nil {
		return err
	}

	// 以下模块都需要认证和权限控制
	// 注册用户管理路由
	registerUserRoutes(r, ctx)

	// 注册角色管理路由
	registerRoleRoutes(r, ctx)

	// 注册 API 权限管理路由
	registerApiPermissionRoutes(r, ctx)

	// 注册组织管理路由
	registerOrgRoutes(r, ctx)

	// 注册菜单管理路由
	registerMenuRoutes(r, ctx)

	// 注册字典管理路由
	registerDictRoutes(r, ctx)

	// 注册配置管理路由
	registerConfigRoutes(r, ctx)

	// 注册登录日志路由
	registerLoginLogRoutes(r, ctx)

	// 注册操作日志路由
	registerOperLogRoutes(r, ctx)

	// 注册附件管理路由
	registerAttachmentRoutes(r, ctx)

	// 注册存储环境管理路由
	registerStorageEnvRoutes(r, ctx)
	return nil
}
