package query

import (
	"fmt"
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

	searchType := c.DefaultQuery("search_type", "text")
	var query string
	if searchType == "text" {
		query = fmt.Sprintf(`PREFIX owl: <http://www.w3.org/2002/07/owl#>
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
prefix pq:<http://nlp_microbe.ccnu.edu.cn/pop/qualifier#>
prefix ps:<http://nlp_microbe.ccnu.edu.cn/pop/statement#>
select distinct ?bacname?disname?reference?bac2name?ref2?bac3name?ref3?bac4name?ref4
where{
    ?bacid rdfs:label %s.
    graph <http://NLPMicrobeKG.com/bac_dis>{
		{?bacid pq:hasAssociation ?dis;
         		pq:isAnnotated ?reference.}
    	?reference pq:hasAnnotation ?dis.
    }
    ?dis rdfs:label ?disname.
    graph <http://NLPMicrobeKG.com/bac_bac>{
        ?bacid pq:hasAssociation ?bac2.
        ?bacid pq:isAnnotated ?ref2.
        ?bac2 pq:isAnnotated ?ref2.
        ?bacid pq:hasNegativeAssociation ?bac3.
        ?bacid pq:isAnnotated ?ref3.
        ?bac3 pq:isAnnotated ?ref3.
        ?bacid pq:hasPositiveAssociation ?bac4.
        ?bacid pq:isAnnotated ?ref4.
        ?bac4 pq:isAnnotated ?ref4.
    }
    ?bac2 rdfs:label ?bac2name.
    ?bac3 rdfs:label ?bac3name.
    ?bac4 rdfs:label ?bac4name.
}`, requestBody.Query)
	} else {
		query = fmt.Sprintf(`PREFIX owl: <http://www.w3.org/2002/07/owl#>
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
prefix pq:<http://nlp_microbe.ccnu.edu.cn/pop/qualifier#>
prefix ps:<http://nlp_microbe.ccnu.edu.cn/pop/statement#>
select distinct ?bacname?modulename?compoundname?mount?unit?foodname
where{
    {?bacid rdfs:label %s;
           pq:hasPathways ?modules;
           pq:hasMetabolites ?compound.}
    ?modules rdfs:label ?modulename.
    ?compound rdfs:label ?compoundname.
    {?food pq:hasNutrients ?compound;
          pq:hasNutrientWeights ?comp_w;
          rdfs:label ?foodname.}
    {?comp_w ps:hasDescriptions ?mount;
             ps:hasMeasuringUnit ?unit.}
    
}`, requestBody.Query)
	}

	res, err := model.QueryInfo(query, requestBody.Limit, requestBody.Offset)
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
