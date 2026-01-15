package v1

import (
	"errors"

	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/request"
	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/evrone/go-clean-template/internal/usecase/user"
	"github.com/gofiber/fiber/v2"
)

// DeductBalance godoc
// @Summary     Deduct user balance
// @Description Deducts specified amount from user balance
// @Tags        balance
// @Accept      json
// @Produce     json
// @Param       request body request.DeductBalance true "Deduct balance request"
// @Success     200 {object} map[string]interface{} "Success response with new balance"
// @Failure     400 {object} response.Error
// @Failure     402 {object} response.Error "Insufficient funds"
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /balance/deduct [post]
func (c *V1) DeductBalance(ctx *fiber.Ctx) error {
	var req request.DeductBalance
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := c.v.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err := c.user.DeductBalance(ctx.Context(), req.UserID, req.Amount)
	if err != nil {
		if errors.Is(err, user.ErrInvalidAmount) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid amount",
			})
		}
		if errors.Is(err, user.ErrUserNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		if errors.Is(err, persistent.ErrInsufficientFunds) {
			return ctx.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
				"error": "insufficient funds",
			})
		}
		c.l.Error(err, "http - v1 - DeductBalance")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	u, err := c.user.GetByID(ctx.Context(), req.UserID)
	if err != nil {
		c.l.Error(err, "http - v1 - DeductBalance - GetByID")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return ctx.JSON(fiber.Map{
		"success":     true,
		"new_balance": u.Balance,
	})
}
