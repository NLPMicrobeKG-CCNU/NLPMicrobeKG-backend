package query

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/handler"
	MDepressionKGQuery "github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/model/MDepressionKG"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/pkg/errno"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/util"
)

// @Summary 查询 MDepression 数据(表显示)
// @Tags search
// @param search_type query string true "查询类型 diseases / food"
// @Param query query string true "查询数据参数"
// @Param limit query int true "LIMIT <= 1000"
// @Param offset query int true "NO RESTRICTIONS"
// @Success 200 {object} []MDepressionKGQuery.DiseaseResponse
// @Success 200 {object} []MDepressionKGQuery.FoodResponse
// @Router /search/mdepression [GET]
func MDepressionQuery(c *gin.Context) {
	var requestBody QueryRequest
	var err error

	queryStr, err := util.ParseBase64(c.DefaultQuery("query", ""))
	if err != nil {
		fmt.Println(err)
		handler.SendError(c, errno.InternalServerError, nil, "parse query failed")
		return
	}

	requestBody.Query = queryStr
	searchType := c.DefaultQuery("search_type", "diseases")
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

	var positiveQuery string
	var negativeQuery string
	var res interface{}

	if searchType == "diseases" {
		positiveQuery = fmt.Sprintf(`PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX entity: <http://nlp_microbe.ccnu.edu.cn/entity#>
prefix pq:<http://nlp_microbe.ccnu.edu.cn/pop/qualifier#>
prefix ps:<http://nlp_microbe.ccnu.edu.cn/pop/statement#>
select distinct ?depression?bac?bacname?syndrome 
where { 
	?depression ps:hasIdentifier 'D003865'.
    ?depression pq:hasPositiveAssociation ?bac.
    {?bac rdf:type entity:BacteriaSpecies;
          rdfs:label ?bacname;
          pq:hasAssociation ?syndrome.}
    {?syndrome rdf:type entity:Disease;
            rdfs:label "%s"@en.}
}`, requestBody.Query)
		negativeQuery = fmt.Sprintf(`PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX entity: <http://nlp_microbe.ccnu.edu.cn/entity#>
prefix pq:<http://nlp_microbe.ccnu.edu.cn/pop/qualifier#>
prefix ps:<http://nlp_microbe.ccnu.edu.cn/pop/statement#>
select distinct ?depression?bac?bacname?syndrome 
where { 
	?depression ps:hasIdentifier 'D003865'.
    ?depression pq:hasNegativeAssociation ?bac.
    {?bac rdf:type entity:BacteriaSpecies;
          rdfs:label ?bacname;
          pq:hasAssociation ?syndrome.}
    {?syndrome rdf:type entity:Disease;
            rdfs:label "%s"@en.}
}`, requestBody.Query)
		res, err = MDepressionKGQuery.DiseaseQuery(positiveQuery, negativeQuery, requestBody.Query, requestBody.Limit, requestBody.Offset)
		if err != nil {
			fmt.Println(err)
			handler.SendError(c, errno.InternalServerError, nil, "query info error")
			return
		}
	} else {
		positiveQuery = fmt.Sprintf(`PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX entity: <http://nlp_microbe.ccnu.edu.cn/entity#>
prefix pq:<http://nlp_microbe.ccnu.edu.cn/pop/qualifier#>
prefix ps:<http://nlp_microbe.ccnu.edu.cn/pop/statement#>
select distinct ?depression?compound?compoundname?food where { 
	?depression ps:hasIdentifier 'D003865'.
    {?depression pq:hasPositiveAssociation  ?compound.}
    {?compound rdf:type entity:Compounds;
            rdfs:label ?compoundname.}
    {?food pq:hasNutrients ?compound;
           rdfs:label "%s".}
}`, requestBody.Query)
		negativeQuery = fmt.Sprintf(`PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
PREFIX entity: <http://nlp_microbe.ccnu.edu.cn/entity#>
prefix pq:<http://nlp_microbe.ccnu.edu.cn/pop/qualifier#>
prefix ps:<http://nlp_microbe.ccnu.edu.cn/pop/statement#>
select distinct ?depression?bac?bacname?compound?compoundname?food where { 
	?depression ps:hasIdentifier 'D003865'.
    {?depression pq:hasNegativeAssociation  ?bac;
                 pq:hasNegativeAssociation  ?compound.}
    {?bac rdf:type entity:BacteriaSpecies;
          rdfs:label ?bacname;
          pq:hasMetabolites ?compound.}
    {?compound rdf:type entity:Compounds;
            rdfs:label ?compoundname.}
    {?food pq:hasNutrients ?compound;
           rdfs:label "%s".}
}`, requestBody.Query)
		res, err = MDepressionKGQuery.FoodQuery(positiveQuery, negativeQuery, requestBody.Query, requestBody.Limit, requestBody.Offset)
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

/*
// @Summary 查询数据(图数据)
// @Tags graph
// @Param search_value query string true "查询数据"
// @Success 200 {object} dbQuery.Data
// @Router /graph [GET]
func MGraphQuery(c *gin.Context) {
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
*/
