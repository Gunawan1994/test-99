package response

import (
	"github.com/labstack/echo/v4"
)

type Response struct {
	Result  bool        `json:"result"`
	Meta    interface{} `json:"meta"`
	Listing interface{} `json:"listings"`
}

func SetResponse(ctx echo.Context, httpstatus int, msg string, meta interface{}, data interface{}, status bool) error {
	return ctx.JSON(httpstatus, Response{
		Result:  status,
		Meta:    meta,
		Listing: data,
	})
}
