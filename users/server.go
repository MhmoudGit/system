package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
	"users/repository"

	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Echo       *echo.Echo
	DB         *pgx.Conn
	Logger     *slog.Logger
	Cfg        *Config
	Ctx        context.Context
	ShutdownCh chan os.Signal
	Server     *http.Server
}

func NewServer() (*Server, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	if err := InitEnv(); err != nil {
		logger.Error("failed to get environment: ", "error", err)
	}
	cfg := LoadConfig()

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, fmt.Sprintf("user=%s password=%s dbname=%s host =%s port=%s", cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbHost, cfg.DbPort))
	if err != nil {
		panic(err)
	}
	logger.Info("database connection established")

	e := echo.New()
	server := &Server{
		Echo:       e,
		Logger:     logger,
		DB:         conn,
		Cfg:        cfg,
		Ctx:        ctx,
		ShutdownCh: make(chan os.Signal, 1),
	}
	signal.Notify(server.ShutdownCh, os.Interrupt)
	return server, nil
}

func (s *Server) Start() {
	s.Server = &http.Server{
		Addr:    s.Cfg.AppAddr,
		Handler: s.Echo,
	}

	go func() {
		if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Error("Failed to listen: ", "error", err)
		}
	}()

	s.Logger.Info("Server running at: " + s.Cfg.AppAddr)
	s.SetupRouter()
	<-s.ShutdownCh
	s.Shutdown()
}

func (s *Server) Shutdown() {
	s.Logger.Info("Initiating shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Graceful shutdown of the server
	if err := s.Server.Shutdown(ctx); err != nil {
		s.Logger.Error("Server forced to shutdown: ", "error", err)
	}

	// Close database connection
	if err := s.DB.Close(ctx); err != nil {
		s.Logger.Error("Failed to disconnect database: ", "error", err)
	}

	s.Logger.Info("Shutdown completed successfully")
}

func (s *Server) SetupRouter() {
	// Middlewares
	s.Echo.Use(middleware.CORSWithConfig(Cors()))
	s.Echo.Use(middleware.Secure())
	s.Echo.Use(middleware.Recover())
	s.Echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				s.Logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				s.Logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
	s.Echo.Validator = &CustomValidator{Validator: validator.New()}

	users := s.Echo.Group("/v1/auth")
	users.GET("/health", func(c echo.Context) error {
		return NewResponse(c, "success", "healthy", "", http.StatusOK)
	})

	auth := AuthHandler{
		DB:     s.DB,
		Repo:   repository.New(s.DB),
		Logger: s.Logger,
		Cfg:    s.Cfg,
		Ctx:    s.Ctx,
	}
	users.POST("/permissions", auth.LoadPermissions)
}

func Cors() middleware.CORSConfig {
	return middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}
}
