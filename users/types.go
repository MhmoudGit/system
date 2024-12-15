package main

import "github.com/labstack/echo/v4"

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Err     string `json:"error,omitempty"`
	Code    int    `json:"code"`
}

func NewResponse(c echo.Context, message string, data any, err string, code int) error {
	res := &Response{
		Message: message,
		Code:    code,
	}

	if err != "" {
		res.Err = err
	}

	if data != nil {
		res.Data = data
	}

	return c.JSON(res.Code, res)
}

func (r Response) Error() string {
	if r.Err != "" {
		return r.Err
	}
	return r.Message
}

type PermissionsDTO struct {
	Permissions []string `json:"permissions" validate:"required"`
}
