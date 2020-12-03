package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	search "github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/handler/query"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/handler/sd"
	"github.com/NLPMicrobeKG-CCNU/NLPMicrobeKG-backend/model"
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

	g.Use(middleware.URLAccessStatistics(model.DB.Redis))

	s := g.Group(apiVersionString + "search")
	s.Use(middleware.IPLimit(model.DB.Redis))
	{
		s.GET("/", search.Query)
	}

	return g
}
