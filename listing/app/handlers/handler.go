package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"listing/app/dto"
	"listing/app/usecases"
	invoice "listing/app/usecases/listing"
	"listing/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	ilog "github.com/sirupsen/logrus"
)

type HTTP struct {
	usecaseListing usecases.Listing
}

func Handler(logger *ilog.Entry, auroradb *sql.DB) *HTTP {
	return &HTTP{
		usecaseListing: invoice.New(logger, auroradb),
	}
}

func (c *HTTP) AddListingHandler(ctx echo.Context) (e error) {
	logger := ctx.Get("logger").(*logrus.Entry)
	logger.Info("handler: AddListingHandler")

	req := dto.Listing{}
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

	e = c.usecaseListing.AddListing(ctx, req)

	return
}

func (c *HTTP) FindAllListingHandler(ctx echo.Context) (e error) {
	logger := ctx.Get("logger").(*logrus.Entry)
	logger.Info("handler: FindAllListingHandler")

	e = c.usecaseListing.FindAllListing(ctx)

	return
}

func (c *HTTP) FindOneListingHandler(ctx echo.Context) (e error) {
	logger := ctx.Get("logger").(*logrus.Entry)
	logger.Info("handler: FindOneListingHandler")

	e = c.usecaseListing.FindOneListing(ctx)

	return
}
