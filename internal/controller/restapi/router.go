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

// NewRouter @title       Skinport API
// @description Skinport items and user balance API
// @version     1.0
// @host        localhost:8080
// @BasePath    /api/v1
func NewRouter(app *fiber.App, cfg *config.Config, l logger.Interface, user usecase.User, items usecase.Items) {
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))

	if cfg.Metrics.Enabled {
		prometheus := fiberprometheus.New("skinport-api")
		prometheus.RegisterAt(app, "/api/metrics")
		app.Use(prometheus.Middleware)
	}

	if cfg.Swagger.Enabled {
		app.Get("/api/swagger/*", swagger.HandlerDefault)
	}

	apiGroup := app.Group("/api")
	apiGroup.Get("/healthz", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })

	apiV1Group := apiGroup.Group("/v1")
	{
		v1.NewRoutes(apiV1Group, l, user, items)
	}

	// Legacy compatibility routes (without /api prefix) to avoid 404s for existing clients.
	app.Get("/healthz", func(ctx *fiber.Ctx) error { return ctx.Redirect("/api/healthz", http.StatusPermanentRedirect) })
	legacyV1Group := app.Group("/v1")
	{
		v1.NewRoutes(legacyV1Group, l, user, items)
	}
}
