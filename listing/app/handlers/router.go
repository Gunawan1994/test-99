package handlers

import (
	"database/sql"
	"net/http"

	"listing/app/middlewares"

	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

type Routes struct {
	AuroraDb *sql.DB
}

func NewRoutes(auroradb *sql.DB) *Routes {
	return &Routes{
		AuroraDb: auroradb,
	}
}

func (route *Routes) RegisterServices(c *echo.Echo) {
	logger := log.WithFields(log.Fields{
		"job":    "RegisterServices",
		"msg_id": xid.New().String(),
	})
	logger.Debug("Running")

	handler := Handler(logger, route.AuroraDb)
	routeListing := c.Group("/listings")

	route.setMiddleware(routeListing)
	routeListing.POST("", handler.AddListingHandler)
	routeListing.GET("", handler.FindAllListingHandler)
	routeListing.GET("/:id", handler.FindOneListingHandler)
}

func (route *Routes) setMiddleware(rGroup *echo.Group) {
	rGroup.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderXRealIP},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost},
	}))

	m := middlewares.New("")
	rGroup.Use(m.AddLoggerToContext, m.DumpRequest)

}
