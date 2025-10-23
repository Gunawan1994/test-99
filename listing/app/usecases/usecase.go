package usecases

import (
	"listing/app/dto"

	"github.com/labstack/echo/v4"
)

type Listing interface {
	AddListing(ctx echo.Context, req dto.Listing) (e error)
	FindAllListing(ctx echo.Context) (e error)
	FindOneListing(ctx echo.Context) (e error)
}
