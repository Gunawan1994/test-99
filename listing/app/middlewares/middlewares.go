package middlewares

import (
	"net/http/httputil"

	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

type middlewaresProcess struct {
	AppName string
}

func New(appName string) middlewaresProcess {
	return middlewaresProcess{
		AppName: appName,
	}
}

func (m middlewaresProcess) AddLoggerToContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		msgID := xid.New().String()

		c.Set("msg_id", msgID)
		c.Set("logger", log.WithFields(log.Fields{
			"app":    m.AppName,
			"msg_id": msgID,
		}))
		return next(c)
	}
}

func (m middlewaresProcess) DumpRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger := c.Get("logger").(*log.Entry)
		requestDump, err := httputil.DumpRequest(c.Request(), true)
		if err != nil {
			logger.WithField("error", err).Error("Catch error")
		}
		logger.Info(string(requestDump))
		return next(c)
	}
}

// GetLogger .
func GetLogger(ctx echo.Context) *log.Entry {
	logger := ctx.Get("logger")
	if logger != nil {
		return logger.(*log.Entry)
	}
	return log.WithFields(log.Fields{
		"app": "assessment",
	})
}
