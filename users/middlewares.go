package main

import (
	"net/http"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

const (
	SuperAdminRole int64 = iota + 1
	AdminRole
	UserRole
)

type JwtCustomClaims struct {
	UserID      int64    `json:"userId"`
	Role        int64    `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

func GenerateToken(secretKey string, userId, role, duration int64, permissions []string) (string, error) {
	claims := &JwtCustomClaims{
		userId,
		role,
		permissions,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(duration))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}

func JWTMiddleware(secretKey string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaims)
		},
		SigningKey: []byte(secretKey),
		ContextKey: "token",
		SuccessHandler: func(c echo.Context) {
			token := c.Get("token").(*jwt.Token)
			claims := token.Claims.(*JwtCustomClaims)
			c.Set("permissions", claims.Permissions)
			c.Set("userID", claims.UserID)
			c.Set("role", claims.Role)
		},
	})

}

func Has(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			permissions, ok := c.Get("permissions").([]string)
			if !ok || !slices.Contains(permissions, permission) {
				err := "invalid permissions"
				return NewResponse(c, "forbidden", nil, err, http.StatusForbidden)
			}

			return next(c)
		}
	}
}