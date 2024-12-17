package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"users/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
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
			if strings.Split(p.Name, ":")[0] == group {
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

// roles handlers

func (h *AuthHandler) GetAllRoles(c echo.Context) error {
	roles, err := h.Repo.ListRoles(h.Ctx)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusInternalServerError)
	}

	if roles == nil {
		roles = make([]repository.Role, 0)
	}

	return NewResponse(c, "success", roles, "", http.StatusOK)
}

func (h *AuthHandler) GetOneRole(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	role, err := h.Repo.GetRole(h.Ctx, int32(id))
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusInternalServerError)
	}

	return NewResponse(c, "success", role, "", http.StatusOK)
}

func (h *AuthHandler) CreateRoles(c echo.Context) error {
	data := new(repository.CreateRoleParams)
	err := c.Bind(data)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	if err = c.Validate(data); err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusUnprocessableEntity)
	}

	role, err := h.Repo.CreateRole(h.Ctx, *data)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	return NewResponse(c, "success", role, "", http.StatusAccepted)
}

func (h *AuthHandler) UpdateRoles(c echo.Context) error {
	data := new(repository.UpdateRoleParams)
	err := c.Bind(data)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	if err = c.Validate(data); err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusUnprocessableEntity)
	}

	err = h.Repo.UpdateRole(h.Ctx, *data)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	return NewResponse(c, "success", nil, "", http.StatusAccepted)
}

func (h *AuthHandler) DeleteRoles(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.Repo.SoftDeleteRole(h.Ctx, int32(id))
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusInternalServerError)
	}

	return NewResponse(c, "success", nil, "", http.StatusOK)
}

// users handlers

func (h *AuthHandler) GetAllUsers(c echo.Context) error {
	users, err := h.Repo.ListUsers(h.Ctx)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusInternalServerError)
	}

	if users == nil {
		users = make([]repository.User, 0)
	}

	usersDTO := make([]UserGetDTO, 0)
	for _, user := range users {
		usersDTO = append(usersDTO, UserGetDTO{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			PhoneNumber: user.PhoneNumber,
			IsActive:    user.IsActive,
			IsVerified:  user.IsVerified,
			Role:        user.Role,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			DeletedAt:   user.DeletedAt,
		})
	}

	return NewResponse(c, "success", usersDTO, "", http.StatusOK)
}

func (h *AuthHandler) GetOneUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.Repo.GetUser(h.Ctx, int32(id))
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusInternalServerError)
	}

	userDTO := UserGetDTO{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		IsActive:    user.IsActive,
		IsVerified:  user.IsVerified,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   user.DeletedAt,
	}

	return NewResponse(c, "success", userDTO, "", http.StatusOK)
}

func (h *AuthHandler) CreateUsers(c echo.Context) error {
	data := new(repository.CreateUserParams)
	err := c.Bind(data)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	if err = c.Validate(data); err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusUnprocessableEntity)
	}

	hashedPassword, err := HashPassword(data.Password)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}
	data.Password = hashedPassword

	user, err := h.Repo.CreateUser(h.Ctx, *data)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	userDTO := UserGetDTO{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		IsActive:    user.IsActive,
		IsVerified:  user.IsVerified,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   user.DeletedAt,
	}

	return NewResponse(c, "success", userDTO, "", http.StatusAccepted)
}

func (h *AuthHandler) UpdateUsers(c echo.Context) error {
	data := new(repository.UpdateUserParams)
	err := c.Bind(data)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	if err = c.Validate(data); err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusUnprocessableEntity)
	}

	if data.Password != "" {
		hashedPassword, err := HashPassword(data.Password)
		if err != nil {
			return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
		}
		data.Password = hashedPassword
	}

	err = h.Repo.UpdateUser(h.Ctx, *data)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	return NewResponse(c, "success", nil, "", http.StatusAccepted)
}

func (h *AuthHandler) DeleteUsers(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.Repo.SoftDeleteUser(h.Ctx, int32(id))
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusInternalServerError)
	}

	return NewResponse(c, "success", nil, "", http.StatusOK)
}

// auth handlers
func VerifyPassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func SendEmail(from, password, to, subject, body string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s", to, subject, body))
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	return err
}

// create a seed for all permissions
func CreatePermissionsSeed(DB *repository.Queries) error {
	permissions := []string{"users:read", "users:create", "users:update", "users:delete", "roles:list", "roles:read", "roles:create", "roles:update", "roles:delete", "permissions:create", "permissions:delete", "permissions:list"}
	for _, permission := range permissions {
		_, err := DB.CreatePermission(context.Background(), permission)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateSuperAdmin(DB *repository.Queries) error {
	role, err := DB.CreateRole(context.Background(), repository.CreateRoleParams{
		RoleName:    "superadmin",
		Permissions: []string{"users:read", "users:create", "users:update", "users:delete", "roles:list", "roles:read", "roles:create", "roles:update", "roles:delete", "permissions:create", "permissions:delete", "permissions:list"},
	})
	if err != nil {
		return err
	}

	hashedPassword, err := HashPassword("superadmin")
	if err != nil {
		return err
	}

	_, err = DB.CreateUser(context.Background(), repository.CreateUserParams{
		Username:    "superadmin",
		Email:       "superadmin@email.com",
		Password:    hashedPassword,
		Role:        int64(role.ID),
		IsActive:    pgtype.Bool{Bool: true, Valid: true},
		IsVerified:  pgtype.Bool{Bool: true, Valid: true},
		FirstName:   pgtype.Text{String: "Super", Valid: true},
		LastName:    pgtype.Text{String: "Admin", Valid: true},
		PhoneNumber: pgtype.Text{String: "1234567890", Valid: true},
	})
	if err != nil {
		return err
	}

	return nil
}

func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return NewResponse(c, "failed", nil, "token is required", http.StatusBadRequest)
	}

	// extract user id from token
	userID := 1
	err := h.Repo.UpdateUser(h.Ctx, repository.UpdateUserParams{
		ID:         int32(userID),
		IsVerified: pgtype.Bool{Bool: true, Valid: true},
	})
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	return NewResponse(c, "success", nil, "", http.StatusOK)
}

func (h *AuthHandler) Register(c echo.Context) error {
	data := new(repository.CreateUserParams)
	err := c.Bind(data)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	if err = c.Validate(data); err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusUnprocessableEntity)
	}

	data.IsActive = pgtype.Bool{Bool: false, Valid: true}
	data.IsVerified = pgtype.Bool{Bool: false, Valid: true}

	hashedPassword, err := HashPassword(data.Password)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}
	data.Password = hashedPassword

	user, err := h.Repo.CreateUser(h.Ctx, *data)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	verificationLink := fmt.Sprintf("%s/verify-email?token=%v", h.Cfg.AppAddr, user.ID)
	verificationBody := fmt.Sprintf("<a href=\"%s\">Verify Email</a>", verificationLink)
	err = SendEmail(h.Cfg.EmailFrom, h.Cfg.EmailPassword, user.Email, "Verify Email", verificationBody)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	userDTO := UserGetDTO{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		IsActive:    user.IsActive,
		IsVerified:  user.IsVerified,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   user.DeletedAt,
	}

	return NewResponse(c, "success", userDTO, "", http.StatusOK)
}

func (h *AuthHandler) Login(c echo.Context) error {
	// take username or email and password
	data := new(LoginDTO)
	err := c.Bind(data)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	if err = c.Validate(data); err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusUnprocessableEntity)
	}

	// check if user exists by user email
	user, err := h.Repo.GetUserByEmail(h.Ctx, data.Email)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	// check if password is correct
	if !VerifyPassword(data.Password, user.Password) {
		return NewResponse(c, "failed", nil, "invalid email or password", http.StatusBadRequest)
	}

	// get user role and permissions
	role, err := h.Repo.GetRole(h.Ctx, int32(user.Role))
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	permissions := role.Permissions

	// generate token
	token, err := GenerateToken(h.Cfg.JWTSecret, int64(user.ID), int64(user.Role), 10, permissions)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	userDTO := UserGetDTO{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		IsActive:    user.IsActive,
		IsVerified:  user.IsVerified,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   user.DeletedAt,
	}

	// return token
	responseData := map[string]interface{}{
		"token": token,
		"user":  userDTO,
	}
	return NewResponse(c, "success", responseData, "", http.StatusOK)
}

// TODO: implement refresh token
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	return NewResponse(c, "success", nil, "", http.StatusOK)
}

func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	data := new(ForgotPasswordDTO)
	err := c.Bind(data)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	if err = c.Validate(data); err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusUnprocessableEntity)
	}

	// generate random password
	password := uuid.New().String()[:9]
	resetPasswordBody := fmt.Sprintf("Your new password is: %s", password)
	err = SendEmail(h.Cfg.EmailFrom, h.Cfg.EmailPassword, data.Email, "Reset Password", resetPasswordBody)
	if err != nil {
		return NewResponse(c, "failed", nil, err.Error(), http.StatusBadRequest)
	}

	return NewResponse(c, "success", nil, "", http.StatusOK)
}

// TODO: implement logout
func (h *AuthHandler) Logout(c echo.Context) error {
	return NewResponse(c, "success", nil, "", http.StatusOK)
}
