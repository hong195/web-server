package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/hong195/web-server/internal/usecase"
	"github.com/hong195/web-server/pkg/logger"
)

func NewRoutes(apiV1Group fiber.Router, l logger.Interface, user usecase.User) {
	c := &V1{
		l:    l,
		v:    validator.New(validator.WithRequiredStructEnabled()),
		user: user,
	}

	apiV1Group.Get("/users/:id", c.GetUser)
	apiV1Group.Post("/balance/deduct", c.DeductBalance)
}
