package restapi

import (
	"net/http"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/hong195/web-server/config"
	_ "github.com/hong195/web-server/docs"
	"github.com/hong195/web-server/internal/controller/restapi/middleware"
	v1 "github.com/hong195/web-server/internal/controller/restapi/v1"
	"github.com/hong195/web-server/internal/usecase"
	"github.com/hong195/web-server/pkg/logger"
)

// @title       Skinport API
// @description Skinport items and user balance API
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(app *fiber.App, cfg *config.Config, l logger.Interface, user usecase.User) {
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))

	if cfg.Metrics.Enabled {
		prometheus := fiberprometheus.New("skinport-api")
		prometheus.RegisterAt(app, "/metrics")
		app.Use(prometheus.Middleware)
	}

	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	app.Get("/healthz", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })

	apiV1Group := app.Group("/v1")
	{
		v1.NewRoutes(apiV1Group, l, user)
	}
}
