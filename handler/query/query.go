package query

import (
	"github.com/gin-gonic/gin"

	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/handler"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/model"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/pkg/errno"
)

type QueryRequest struct {
	Query  string `json:"query"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

// @Summary 查询数据
// @Tags search
// @Param data body query.QueryRequest true "查询数据参数"
// @Success 200 "OK"
// @Router /search [GET]
func Query(c *gin.Context) {
	var requestBody QueryRequest
	if err := c.Bind(&requestBody); err != nil {
		handler.SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	res, err := model.QueryInfo(requestBody.Query, requestBody.Limit, requestBody.Offset)
	if err != nil {
		handler.SendError(c, errno.InternalServerError, nil, "query info error")
		return
	}

	handler.SendResponse(c, nil, handler.Response{
		Code:    200,
		Message: "ok",
		Data:    res,
	})
}
