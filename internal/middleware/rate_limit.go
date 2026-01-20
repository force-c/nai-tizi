package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/force-c/nai-tizi/internal/domain/response"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// 简易限流：按用户ID/IP在固定窗口内计数，超过限制返回429
func RateLimit(rc *redis.Client, keyPrefix string, limit int, windowSeconds int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		v, exists := c.Get("userId")
		uid := ""
		if exists {
			uid = fmt.Sprint(v)
		}
		if uid == "" {
			uid = c.ClientIP()
		}
		key := keyPrefix + ":" + uid
		n, err := rc.Incr(ctx, key).Result()
		if err != nil {
			c.AbortWithStatusJSON(200, gin.H{"code": response.CodeTooManyRequests, "msg": "限流中"})
			return
		}
		if n == 1 {
			rc.Expire(ctx, key, time.Duration(windowSeconds)*time.Second)
		}
		if n > int64(limit) {
			c.AbortWithStatusJSON(200, gin.H{"code": response.CodeTooManyRequests, "msg": "请求过于频繁"})
			return
		}
		c.Next()
	}
}
