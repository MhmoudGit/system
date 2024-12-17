package main

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
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

// permissions handlers

func (h *AuthHandler) GetAllPermissions(c echo.Context) error {
	permissions, err := h.Repo.ListPermissions(h.Ctx)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusInternalServerError)
	}

	if permissions == nil {
		permissions = make([]repository.Permission, 0)
	}

	if c.QueryParam("group") != "" {
		group := c.QueryParam("group")
		filteredPermissions := make([]repository.Permission, 0)
		for _, p := range permissions {
			if strings.Split(p.Name, ".")[0] == group {
				filteredPermissions = append(filteredPermissions, p)
			}
		}
		permissions = filteredPermissions
	}

	return NewResponse(c, "success", permissions, "", http.StatusOK)
}

func (h *AuthHandler) CreatePermissions(c echo.Context) error {
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
		if err != nil {
			return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
		}
	}

	if err := tx.Commit(h.Ctx); err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	return NewResponse(c, "success", data, "", http.StatusAccepted)
}

func (h *AuthHandler) DeletePermissions(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.Repo.SoftDeletePermission(h.Ctx, int32(id))
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusInternalServerError)
	}

	return NewResponse(c, "success", nil, "", http.StatusOK)
}
