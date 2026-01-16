package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hong195/web-server/internal/controller/restapi/v1/response"
)

// GetItems godoc
// @Summary     List Skinport items
// @Description Returns Skinport items with tradable and non-tradable minimum prices
// @Tags        items
// @Produce     json
// @Success     200 {array} response.ItemResponse
// @Failure     500 {object} response.Error "internal server error"
// @Failure     502 {object} response.Error "failed to fetch items from skinport"
// @Router      /items [get]
func (c *V1) getItems(ctx *fiber.Ctx) error {
	if c.items == nil {
		c.l.Error("items usecase is not configured")
		return errorResponse(ctx, fiber.StatusInternalServerError, "internal server error")
	}

	items, err := c.items.GetItems(ctx.Context())
	if err != nil {
		c.l.Error(err, "http - v1 - getItems")
		return errorResponse(ctx, fiber.StatusBadGateway, "failed to fetch items from skinport")
	}

	resp := make([]response.ItemResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, response.ItemResponse{
			MarketHashName:      item.MarketHashName,
			TradableMinPrice:    item.MinPriceTradable,
			NonTradableMinPrice: item.MinPriceNonTradable,
		})
	}

	return ctx.JSON(resp)
}
