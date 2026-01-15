package v1

import (
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// NewRoutes -.
func NewRoutes(apiV1Group fiber.Router, l logger.Interface) {
	_ = &V1{l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	// Routes will be added here
	_ = apiV1Group
}
