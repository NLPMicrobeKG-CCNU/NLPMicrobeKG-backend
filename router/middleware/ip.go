package middleware

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"

	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/handler"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/pkg/errno"
)

// IPLimit 限制ip访问次数
func IPLimit(pool *redis.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn := pool.Get()
		defer conn.Close()

		// 获取IP地址
		ip := c.ClientIP()
		// 检查IP地址是否在redis中
		value, err := conn.Do("GET", ip)
		if err != nil {
			handler.SendError(c, err, nil, "redis get ip error")
			c.Abort()
			return
		}

		// 如果在，直接返回，因为他快速调用了多次redis
		if value != nil {
			handler.SendResponse(c, errno.ErrRequestTooQuick, nil)
			c.Abort()
			return
		}

		// 不在，则设置一个，之后直接next
		_, err = conn.Do("SET", ip, "placeholder", "EX", "10")
		if err != nil {
			handler.SendError(c, err, nil, "redis set ip error")
			c.Abort()
			return
		}
		c.Next()
	}
}
