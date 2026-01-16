package v1

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hong195/web-server/internal/usecase/user"
)

// GetUser godoc
// @Summary     Get user by ID
// @Description Returns user with their balance
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Success     200 {object} entity.User
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /users/{id} [get]
func (c *V1) GetUser(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid user id")
	}

	u, err := c.user.GetByID(ctx.Context(), userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return errorResponse(ctx, fiber.StatusNotFound, "user not found")
		}
		c.l.Error(err, "http - v1 - GetUser")
		return errorResponse(ctx, fiber.StatusInternalServerError, "internal server error")
	}

	return ctx.JSON(u)
}
