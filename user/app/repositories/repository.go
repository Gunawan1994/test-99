package repositories

import (
	"user/app/entity"

	"github.com/labstack/echo/v4"
)

type User interface {
	CreateUser(c echo.Context, paramsListing map[string]interface{}) (id int64, e error)
	GetAllUser(c echo.Context, meta entity.Meta) (data []entity.User, total int64, e error)
	GetUser(c echo.Context, id int64) (data entity.User, e error)
}
