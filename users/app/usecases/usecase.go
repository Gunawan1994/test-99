package usecases

import (
	"user/app/dto"

	"github.com/labstack/echo/v4"
)

type User interface {
	AddUser(ctx echo.Context, req dto.User) (e error)
	FindAllUser(ctx echo.Context) (e error)
	FindOneUser(ctx echo.Context) (e error)
}
