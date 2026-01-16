package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hong195/web-server/internal/controller/restapi/v1/response"
)

func errorResponse(ctx *fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(response.Error{Error: msg})
}
