package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hong195/web-server/internal/controller/restapi/v1/response"
)

const defaultItemsLimit = 100

// GetItems godoc
// @Summary     List Skinport items
// @Description Returns Skinport items with tradable and non-tradable minimum prices (paginated)
// @Tags        items
// @Produce     json
// @Param       page  query int false "Page number" default(1)
// @Param       limit query int false "Items per page" default(100)
// @Success     200 {object} response.ItemsPagedResponse
// @Failure     500 {object} response.Error "internal server error"
// @Failure     502 {object} response.Error "failed to fetch items from skinport"
// @Router      /items [get]
func (c *V1) getItems(ctx *fiber.Ctx) error {
	if c.items == nil {
		c.l.Error("items usecase is not configured")
		return errorResponse(ctx, fiber.StatusInternalServerError, "internal server error")
	}

	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", defaultItemsLimit)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > defaultItemsLimit {
		limit = defaultItemsLimit
	}

	items, err := c.items.GetItems(ctx.Context())
	if err != nil {
		c.l.Error(err, "http - v1 - getItems")
		return errorResponse(ctx, fiber.StatusBadGateway, "failed to fetch items from skinport")
	}

	total := len(items)
	totalPages := (total + limit - 1) / limit

	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	pagedItems := items[start:end]

	resp := make([]response.ItemResponse, 0, len(pagedItems))
	for _, item := range pagedItems {
		resp = append(resp, response.ItemResponse{
			MarketHashName:      item.MarketHashName,
			TradableMinPrice:    item.MinPriceTradable,
			NonTradableMinPrice: item.MinPriceNonTradable,
		})
	}

	return ctx.JSON(response.ItemsPagedResponse{
		Items:      resp,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}
