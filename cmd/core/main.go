package main

import (
	sensor "andreas/internal/handler"
	"andreas/internal/repository/mysql"
	"andreas/internal/service/core"
	"context"
	"embed"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	_ "github.com/swaggo/echo-swagger/example/docs"
)

type Dependency struct {
	Handler struct {
		Sensor sensor.Dependency
	}
	Service struct {
		Core core.Dependency
	}
}

type Configuration struct {
	Repository struct {
		Mysql mysql.Configuration
	}
	Service struct {
		Core core.Configuration
	}
}

//go:embed files/*
var env embed.FS

// LoadEnv doing process split string from file .env
// and extract each key and value to os environment
func LoadEnv(env string) {
	s := strings.Split(env, "\n")

	for _, v := range s {
		if len(v) == 0 || !strings.Contains(v, "=") {
			continue
		}

		vS := strings.SplitN(v, "=", 2)
		os.Setenv(vS[0], vS[1])
	}
}

func initialize(log zerolog.Logger) {
	key_value, err := env.ReadFile("files/dev.env")
	if err != nil {
		log.Error().Msg("failed initialize env: " + err.Error())
	}
	LoadEnv(string(key_value))
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server for api.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v2
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	logger := zerolog.New(os.Stdout)
	defer cancel()
	initialize(logger)

	d := new(Dependency)
	c := new(Configuration)

	{
		// initialize configuration
		c.Repository.Mysql = mysql.Configuration{
			DSN: os.Getenv("DSN"),
		}
		c.Service.Core = core.Configuration{
			Port:        8080,
			ReadTimeout: 30 * time.Second,
		}
	}

	db := mysql.New(c.Repository.Mysql, mysql.Dependency{
		Logger: logger,
	})

	{
		// initialize depedency
		d.Handler.Sensor = sensor.Dependency{
			MySQL:  db,
			Logger: logger,
		}
		d.Service.Core = core.Dependency{
			Sensor: sensor.New(d.Handler.Sensor),
			Logger: logger,
		}
	}

	core := core.New(c.Service.Core, d.Service.Core)

	if err := core.Serve(); err != nil {
		logger.Fatal().Msg("Failed to start " + err.Error())
	}

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := core.Shutdown(ctx); err != nil {
		logger.Fatal().Msg("Shutdown indopass service")
	}
}
