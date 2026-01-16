package entity

type Item struct {
	MarketHashName      string   `json:"market_hash_name"`
	Currency            string   `json:"currency"`
	SuggestedPrice      *float64 `json:"suggested_price"`
	ItemPage            string   `json:"item_page"`
	MarketPage          string   `json:"market_page"`
	MinPriceTradable    *float64 `json:"min_price_tradable"`
	MinPriceNonTradable *float64 `json:"min_price_non_tradable"`
}
