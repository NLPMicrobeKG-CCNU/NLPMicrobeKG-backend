package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/docs"
	search "github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/handler/query"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/handler/sd"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/router/middleware"
)

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// The health check handlers
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	apiVersionString := "/api/v1"

	//g.Use(middleware.URLAccessStatistics(model.DB.Redis))

	s := g.Group(apiVersionString + "/search")
	// search for default search
	// type: text / data
	//s.Use(middleware.IPLimit(model.DB.Redis))
	{
		// /api/v1/search/
		s.GET("", search.Query)
	}

	graph := g.Group(apiVersionString + "/graph")
	// search for graph query
	// graph search
	//graph.Use(middleware.IPLimit(model.DB.Redis))
	{
		// /api/v1/graph
		graph.GET("", search.GraphQuery)
	}

	// swagger API doc
	swaggerRouter := g.Group(apiVersionString + "/swagger/*any")
	{
		swaggerRouter.GET("", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return g
}
