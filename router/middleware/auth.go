package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/handler"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/pkg/errno"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/pkg/token"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the json web token.
		if _, err := token.ParseRequest(c); err != nil {
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
