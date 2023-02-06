package core

import (
	sensor "andreas/internal/handler"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type CoreService interface {
	Serve() error
	Shutdown(ctx context.Context) error
}

type Configuration struct {
	Port              int
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
}

type Dependency struct {
	Sensor sensor.Sensor
	zerolog.Logger
}

// validator request
type CoreValidator struct {
	validator *validator.Validate
}

func (cv *CoreValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func New(cfg Configuration, dep Dependency) CoreService {
	handler := echo.New()
	handler.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  log.ERROR,
	}))

	handler.Validator = &CoreValidator{validator: validator.New()}

	handler.GET("/data", dep.Sensor.GetData)
	handler.POST("/data", dep.Sensor.StoreData)
	handler.GET("/swagger/*", echoSwagger.WrapHandler)

	srv := http.Server{
		Addr:        fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout: cfg.ReadTimeout,
		Handler:     handler,
	}

	return &instance{cfg, dep, &srv}
}

type instance struct {
	Configuration
	Dependency

	handler *http.Server
}

func (x *instance) Serve() error {
	x.Logger.Info().Msg(fmt.Sprintf("Service started at port %d...", x.Configuration.Port))
	return x.handler.ListenAndServe()
}

func (x *instance) Shutdown(ctx context.Context) error {
	return x.handler.Shutdown(ctx)
}
