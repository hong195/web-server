package v1

import (
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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
