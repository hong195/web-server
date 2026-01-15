package v1

import (
	"errors"
	"strconv"

	"github.com/evrone/go-clean-template/internal/usecase/user"
	"github.com/gofiber/fiber/v2"
)

// GetUser godoc
// @Summary     Get user by ID
// @Description Returns user with their balance
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Success     200 {object} entity.User
// @Failure     400 {object} map[string]string "Invalid user ID"
// @Failure     404 {object} map[string]string "User not found"
// @Failure     500 {object} map[string]string "Internal server error"
// @Router      /users/{id} [get]
func (c *V1) GetUser(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	u, err := c.user.GetByID(ctx.Context(), userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		c.l.Error(err, "http - v1 - GetUser")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return ctx.JSON(u)
}
