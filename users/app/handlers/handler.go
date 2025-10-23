package handlers

import (
	"fmt"
	"net/http"

	"user/app/dto"
	"user/app/usecases"
	user "user/app/usecases/user"
	"user/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	ilog "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type HTTP struct {
	usecaseUser usecases.User
}

func Handler(logger *ilog.Entry, db *mongo.Client) *HTTP {
	return &HTTP{
		usecaseUser: user.New(logger, db),
	}
}

func (c *HTTP) AddUserHandler(ctx echo.Context) (e error) {
	logger := ctx.Get("logger").(*logrus.Entry)
	logger.Info("handler: AddUserHandler")

	req := dto.User{}
	if e = ctx.Bind(&req); e != nil {
		logger.WithField("error", e.Error()).Error("Catch error bind request")
		e = response.SetResponse(ctx, http.StatusBadRequest,
			"Missing mandatory parameter", nil, nil, false)
		return
	}

	validate := validator.New()
	if e = validate.Struct(&req); e != nil {
		errs := e.(validator.ValidationErrors)
		for _, fieldErr := range errs {
			logger.WithField("error", e.Error()).Error(fmt.Printf("field %s: %s\n", fieldErr.Field(), fieldErr.Tag()))
			e = response.SetResponse(ctx, http.StatusBadRequest,
				fmt.Sprintf("Missing mandatory parameter %s", fieldErr.Field()), nil, nil, false)
			return
		}
		return
	}

	e = c.usecaseUser.AddUser(ctx, req)

	return
}

func (c *HTTP) FindAllUserHandler(ctx echo.Context) (e error) {
	logger := ctx.Get("logger").(*logrus.Entry)
	logger.Info("handler: FindAllUserHandler")

	e = c.usecaseUser.FindAllUser(ctx)

	return
}

func (c *HTTP) FindOneUserHandler(ctx echo.Context) (e error) {
	logger := ctx.Get("logger").(*logrus.Entry)
	logger.Info("handler: FindOneUserHandler")

	e = c.usecaseUser.FindOneUser(ctx)

	return
}
