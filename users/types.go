package main

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

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

type UserGetDTO struct {
	ID          int32            `json:"id"`
	Username    string           `json:"username"`
	Email       string           `json:"email"`
	FirstName   pgtype.Text      `json:"first_name"`
	LastName    pgtype.Text      `json:"last_name"`
	PhoneNumber pgtype.Text      `json:"phone_number"`
	IsActive    pgtype.Bool      `json:"is_active"`
	IsVerified  pgtype.Bool      `json:"is_verified"`
	Role        int64            `json:"role"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
	DeletedAt   pgtype.Timestamp `json:"deleted_at"`
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ForgotPasswordDTO struct {
	Email string `json:"email" validate:"required"`
}
