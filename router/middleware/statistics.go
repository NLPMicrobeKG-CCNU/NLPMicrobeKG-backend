package middleware

import (
	"regexp"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"

	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/handler"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/model"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/util"
)

// URLAccessStatistics 记录api调用次数
func URLAccessStatistics(conn *redis.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		redisConn := model.DB.Redis.Get()
		defer redisConn.Close()
		// 获取URL地址
		url := c.FullPath()

		// Skip for swagger files
		reg := regexp.MustCompile("swagger")
		if reg.MatchString(url) {
			return
		}

		// 去除无用信息
		url = strings.Replace(url, "/api/v1/", "", 1)

		// 新增记录
		err := addNewURLAccessRecord(url, &redisConn)
		if err != nil {
			handler.SendError(c, err, nil, "redis get url error")
			c.Abort()
			return
		}
		c.Next()
	}
}

func addNewURLAccessRecord(url string, conn *redis.Conn) error {
	now := util.GetCurrentDate()

	// 添加记录
	_, err := (*conn).Do("zincrby", now+"URLRecord", 1, url)
	if err != nil {
		return err
	}

	return nil
}
