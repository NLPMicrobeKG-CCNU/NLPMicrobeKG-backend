package handler

import (
	"net/http"

	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/log"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/pkg/errno"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	code, message := errno.DecodeErr(err)
	log.Info(message,
		zap.String("X-Request-Id", util.GetReqID(c)))

	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func SendBadRequest(c *gin.Context, err error, data interface{}, cause string) {
	code, message := errno.DecodeErr(err)
	log.Error(message,
		zap.String("X-Request-Id", util.GetReqID(c)),
		zap.String("cause", cause))

	c.JSON(http.StatusBadRequest, Response{
		Code:    code,
		Message: message + ": " + cause,
		Data:    data,
	})
}

func SendError(c *gin.Context, err error, data interface{}, cause string) {
	code, message := errno.DecodeErr(err)
	log.Error(message,
		zap.String("X-Request-Id", util.GetReqID(c)),
		zap.String("cause", cause))

	c.JSON(http.StatusInternalServerError, Response{
		Code:    code,
		Message: message + ": " + cause,
		Data:    data,
	})
}
