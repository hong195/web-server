package response

// ItemResponse represents the public response for a Skinport item.
type ItemResponse struct {
	MarketHashName      string   `json:"market_hash_name"`
	TradableMinPrice    *float64 `json:"tradable_min_price"`
	NonTradableMinPrice *float64 `json:"non_tradable_min_price"`
}
