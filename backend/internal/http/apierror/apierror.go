package apierror

import "github.com/labstack/echo/v4"

type Response struct {
	Error Body `json:"error"`
}

type Body struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func JSON(c echo.Context, status int, code string, message string) error {
	return c.JSON(status, Response{
		Error: Body{
			Code:    code,
			Message: message,
		},
	})
}
