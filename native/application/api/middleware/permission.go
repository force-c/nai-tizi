package middleware

import (
	"github.com/gcc798/lightning/application/api/domain/response"
	"github.com/gcc798/lightning/application/api/service"
	"github.com/labstack/echo/v5"
)

// Permission checks one API permission after authentication.
func Permission(permissionService service.PermissionService, resource string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userID, ok := contextUserID(c)
			if !ok {
				response.Forbidden(c, "用户信息不存在")
				return nil
			}
			allowed, err := permissionService.CheckPermission(c.Request().Context(), userID, resource, resourceAction(resource))
			if err != nil {
				response.InternalServerError(c, "权限检查失败: "+err.Error())
				return nil
			}
			if !allowed {
				response.Forbidden(c, "无权限访问")
				return nil
			}
			return next(c)
		}
	}
}

// PermissionAny allows the request when any resource is granted.
func PermissionAny(permissionService service.PermissionService, resources []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userID, ok := contextUserID(c)
			if !ok {
				response.Forbidden(c, "用户信息不存在")
				return nil
			}
			for _, resource := range resources {
				allowed, err := permissionService.CheckPermission(c.Request().Context(), userID, resource, resourceAction(resource))
				if err == nil && allowed {
					return next(c)
				}
			}
			response.Forbidden(c, "无权限访问")
			return nil
		}
	}
}

// PermissionAll requires every resource to be granted.
func PermissionAll(permissionService service.PermissionService, resources []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userID, ok := contextUserID(c)
			if !ok {
				response.Forbidden(c, "用户信息不存在")
				return nil
			}
			for _, resource := range resources {
				allowed, err := permissionService.CheckPermission(c.Request().Context(), userID, resource, resourceAction(resource))
				if err != nil || !allowed {
					response.Forbidden(c, "无权限访问")
					return nil
				}
			}
			return next(c)
		}
	}
}

func contextUserID(c *echo.Context) (int64, bool) {
	value := c.Get("userId")
	userID, ok := value.(int64)
	return userID, ok
}

func resourceAction(resource string) string {
	if len(resource) > 5 && resource[len(resource)-5:] == ".read" {
		return "read"
	}
	return "write"
}
