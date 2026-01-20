package middleware

import (
	"strings"

	"github.com/force-c/nai-tizi/internal/config"
	"github.com/force-c/nai-tizi/internal/domain/model"
	"github.com/force-c/nai-tizi/internal/domain/response"
	"github.com/force-c/nai-tizi/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Auth 认证中间件
// 1. 从配置的请求头读取 AccessToken
// 2. 验证 AccessToken
// 3. 查询用户的组织ID
// 4. 设置用户信息到 context
func Auth(tokenManager service.TokenManager, cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从配置的请求头读取 Token
		tokenHeader := cfg.Auth.TokenHeader
		token := c.GetHeader(tokenHeader)
		if token == "" {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 去除 "Bearer " 前缀
		token = strings.TrimPrefix(token, "Bearer ")

		// 验证 AccessToken
		claims, err := tokenManager.ValidateAccessToken(c.Request.Context(), token)
		if err != nil {
			response.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		// 可选：验证请求头中的 clientId 是否与 Token 中的一致
		headerClientId := c.GetHeader("clientid")
		if headerClientId == "" {
			headerClientId = c.Query("clientid")
		}
		if headerClientId != "" && claims.ClientId != headerClientId {
			response.Unauthorized(c, "客户端ID与Token不匹配")
			c.Abort()
			return
		}

		// 查询用户的组织ID
		var user model.User
		if err := db.Select("org_id").Where("id = ?", claims.UserId).First(&user).Error; err == nil {
			// 设置用户的组织ID到 context
			c.Set("orgId", user.OrgId)
		}

		// 设置用户信息到 context
		c.Set("userId", claims.UserId)
		c.Set("userName", claims.UserName)
		c.Set("clientId", claims.ClientId)
		c.Set("deviceType", claims.DeviceType)
		c.Next()
	}
}
