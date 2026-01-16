package response

// ItemResponse represents the public response for a Skinport item.
type ItemResponse struct {
	MarketHashName      string   `json:"market_hash_name"`
	TradableMinPrice    *float64 `json:"tradable_min_price"`
	NonTradableMinPrice *float64 `json:"non_tradable_min_price"`
}

// ItemsPagedResponse represents a paginated list of items.
type ItemsPagedResponse struct {
	Items      []ItemResponse `json:"items"`
	Page       int            `json:"page" example:"1"`
	Limit      int            `json:"limit" example:"100"`
	Total      int            `json:"total" example:"5000"`
	TotalPages int            `json:"total_pages" example:"50"`
}
