package main

import (
	"database/sql"
	"net/http"
	"os"

	"listing/app/handlers"
	"listing/pkg/db"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func initLogger() {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.999",
	})
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Warn(".env file not found")
	}
}

func configAndStartServer() {
	initLogger()
	loadEnv()
	dbClient := dbClient()

	htmlEcho := setWebRouter(dbClient)
	start(htmlEcho)
}

func setWebRouter(dbClient *sql.DB) *echo.Echo {
	e := echo.New()

	root := e.Group("")
	root.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "PONG!")
	})

	// Register Routes
	handlers.NewRoutes(dbClient).RegisterServices(e)

	return e
}

func start(htmlEcho *echo.Echo) {
	// Start Run HTML Echo
	if err := htmlEcho.Start(os.Getenv("LISTEN_PORT")); err != nil {
		log.WithField("error", err).Error("Unable to start the server")
		os.Exit(1)
	}
}

func dbClient() *sql.DB {
	return db.NewConn()
}

func main() {
	initLogger()
	configAndStartServer()
}
