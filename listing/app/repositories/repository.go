package repositories

import (
	"listing/app/entity"

	"github.com/labstack/echo/v4"
)

type Listing interface {
	CreateListing(c echo.Context, paramsListing map[string]interface{}) (id int64, e error)
	GetAllListing(c echo.Context, meta entity.Meta) (data []entity.Listing, total int64, e error)
	GetListing(c echo.Context, id string) (data entity.Listing, e error)
}
