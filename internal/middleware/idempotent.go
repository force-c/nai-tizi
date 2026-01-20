package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/force-c/nai-tizi/internal/domain/response"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	idempotentHeader = "Idempotency-Token"
	defaultTTL       = 10 * time.Minute
	keyPrefix        = "idempotent:"
)

// Idempotent 幂等性中间件
func Idempotent(rc *redis.Client, ttl time.Duration) gin.HandlerFunc {
	if ttl <= 0 {
		ttl = defaultTTL
	}

	return func(c *gin.Context) {
		token := c.GetHeader(idempotentHeader)
		if token == "" {
			response.FailCode(c, response.CodeBadRequest, "缺少 Idempotency-Token")
			c.Abort()
			return
		}

		ctx := context.Background()
		key := keyPrefix + token

		ok, err := rc.SetNX(ctx, key, c.FullPath(), ttl).Result()
		if err != nil {
			response.FailCode(c, response.CodeServerError, "幂等校验失败")
			c.Abort()
			return
		}

		if !ok {
			response.FailCode(c, response.CodeBadRequest, fmt.Sprintf("请求已处理: %s", token))
			c.Abort()
			return
		}

		c.Next()
	}
}
