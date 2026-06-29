package handlers

import "github.com/labstack/echo/v4"

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func JSONError(c echo.Context, status int, code string, message string) error {
	return c.JSON(status, ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}
