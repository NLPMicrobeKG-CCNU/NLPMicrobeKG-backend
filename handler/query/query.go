package query

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/handler"
	dbQuery "github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/model/query"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/pkg/errno"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/util"
)

type QueryRequest struct {
	Query  string `json:"query"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

// @Summary 查询数据(表显示)
// @Tags search
// @param search_type query string true "查询类型 text / data"
// @Param query query string true "查询数据参数"
// @Param limit query int true "LIMIT <= 1000"
// @Param offset query int true "NO RESTRICTIONS"
// @Success 200 {object} []dbQuery.TextResponse
// @Success 200 {object} []dbQuery.DataResponse
// @Router /search [GET]
func MicrobeKGQuery(c *gin.Context) {
	var requestBody QueryRequest
	var err error

	requestBody.Query = util.FormatRequestQuery(c.DefaultQuery("query", ""))
	searchType := c.DefaultQuery("search_type", "text")
	limitStr := c.DefaultQuery("limit", "1000")
	pageStr := c.DefaultQuery("page", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		fmt.Println(err)
		handler.SendError(c, errno.InternalServerError, nil, "query info error")
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		fmt.Println(err)
		handler.SendError(c, errno.InternalServerError, nil, "query info error")
		return
	}

	if limit > 1000 {
		limit = 1000
	}
	requestBody.Limit = limit
	requestBody.Offset = page * limit

	var query string
	var res interface{}

	if searchType == "text" {
		query = fmt.Sprintf(`PREFIX owl: <http://www.w3.org/2002/07/owl#>
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
prefix pq:<http://nlp_microbe.ccnu.edu.cn/pop/qualifier#>
prefix ps:<http://nlp_microbe.ccnu.edu.cn/pop/statement#>
select distinct ?bacname?disname?reference?bac2name?ref2?bac3name?ref3?bac4name?ref4
where{
    ?bacid rdfs:label "%s".
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
		res, err = dbQuery.TextQuery(query, requestBody.Query, requestBody.Limit, requestBody.Offset)
		if err != nil {
			fmt.Println(err)
			handler.SendError(c, errno.InternalServerError, nil, "query info error")
			return
		}
	} else {
		query = fmt.Sprintf(`PREFIX owl: <http://www.w3.org/2002/07/owl#>
PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
prefix pq:<http://nlp_microbe.ccnu.edu.cn/pop/qualifier#>
prefix ps:<http://nlp_microbe.ccnu.edu.cn/pop/statement#>
select distinct ?bacname?modulename?compoundname?mount?unit?foodid?foodname
where{
    {?bacid rdfs:label "%s";
           pq:hasPathways ?modules;
           pq:hasMetabolites ?compound.}
    ?modules rdfs:label ?modulename.
    ?compound rdfs:label ?compoundname.
    {?food pq:hasNutrients ?compound;
          pq:hasNutrientWeights ?comp_w;
          ps:hasIdentifier ?foodid;
          rdfs:label ?foodname.}
    {?comp_w ps:hasDescriptions ?mount;
             ps:hasMeasuringUnit ?unit.}
    
}`, requestBody.Query)
		res, err = dbQuery.DataQuery(query, requestBody.Query, requestBody.Limit, requestBody.Offset)
		if err != nil {
			handler.SendError(c, errno.InternalServerError, nil, "query info error")
			return
		}
	}

	handler.SendResponse(c, nil, handler.Response{
		Code:    200,
		Message: "ok",
		Data:    res,
	})
}

// @Summary 查询数据(图数据)
// @Tags graph
// @Param search_value query string true "查询数据"
// @Success 200 {object} dbQuery.Data
// @Router /graph [GET]
func GraphQuery(c *gin.Context) {
	searchValue := c.DefaultQuery("search_value", "")

	query, err := dbQuery.ConvertQueryString(searchValue)
	if err != nil {
		fmt.Println(err)
		handler.SendError(c, errno.InternalServerError, nil, "query info error")
		return
	}

	res, err := dbQuery.GraphQuery(query)
	if err != nil {
		fmt.Println(err)
		handler.SendError(c, errno.InternalServerError, nil, "query info error")
		return
	}

	handler.SendResponse(c, nil, handler.Response{
		Code:    200,
		Message: "ok",
		Data:    res,
	})
}
