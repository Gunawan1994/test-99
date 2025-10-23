package handlers

import (
	"net/http"

	"user/app/middlewares"

	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

type Routes struct {
	Db *mongo.Client
}

func NewRoutes(db *mongo.Client) *Routes {
	return &Routes{
		Db: db,
	}
}

func (route *Routes) RegisterServices(c *echo.Echo) {
	logger := log.WithFields(log.Fields{
		"job":    "RegisterServices",
		"msg_id": xid.New().String(),
	})
	logger.Debug("Running")

	handler := Handler(logger, route.Db)
	routeListing := c.Group("/users")

	route.setMiddleware(routeListing)
	routeListing.POST("", handler.AddUserHandler)
	routeListing.GET("", handler.FindAllUserHandler)
	routeListing.GET("/:id", handler.FindOneUserHandler)
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
