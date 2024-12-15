package main

import (
	"context"
	"log/slog"
	"net/http"
	"users/repository"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	DB     *pgx.Conn
	Repo   *repository.Queries
	Logger *slog.Logger
	Cfg    *Config
	Ctx    context.Context
}

func (h *AuthHandler) LoadPermissions(c echo.Context) error {
	data := new(PermissionsDTO)
	err := c.Bind(data)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	if err = c.Validate(data); err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusUnprocessableEntity)
	}

	tx, err := h.DB.Begin(h.Ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(h.Ctx)

	qtx := h.Repo.WithTx(tx)
	for _, p := range data.Permissions {
		_, err := qtx.CreatePermission(h.Ctx, p)
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	if err := tx.Commit(h.Ctx); err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	return NewResponse(c, "success", data, "", http.StatusAccepted)
}
