package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/andybalholm/brotli"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/hong195/web-server/config"
	"github.com/hong195/web-server/internal/entity"
)

// skinportItem represents the API response structure from Skinport.
type skinportItem struct {
	MarketHashName string   `json:"market_hash_name"`
	Currency       string   `json:"currency"`
	SuggestedPrice *float64 `json:"suggested_price"`
	ItemPage       string   `json:"item_page"`
	MarketPage     string   `json:"market_page"`
	MinPrice       *float64 `json:"min_price"`
	MaxPrice       *float64 `json:"max_price"`
	MeanPrice      *float64 `json:"mean_price"`
	MedianPrice    *float64 `json:"median_price"`
	Quantity       int      `json:"quantity"`
	CreatedAt      int64    `json:"created_at"`
	UpdatedAt      int64    `json:"updated_at"`
}

// SkinportRepo implements repo.ItemsRepo using Skinport HTTP API.
type SkinportRepo struct {
	client   *http.Client
	baseURL  string
	appID    int
	currency string
}

// NewSkinportRepo creates a new SkinportRepo.
func NewSkinportRepo(client *http.Client, cfg config.Skinport) *SkinportRepo {
	return &SkinportRepo{
		client:   client,
		baseURL:  cfg.BaseURL,
		appID:    cfg.AppID,
		currency: cfg.Currency,
	}
}

// GetItems fetches items from Skinport API and merges tradable/non-tradable prices.
func (r *SkinportRepo) GetItems(ctx context.Context) ([]entity.Item, error) {
	// Fetch tradable and non-tradable items in parallel
	type result struct {
		items []skinportItem
		err   error
	}

	tradableCh := make(chan result, 1)
	nonTradableCh := make(chan result, 1)

	go func() {
		items, err := r.fetchItems(ctx, true)
		tradableCh <- result{items: items, err: err}
	}()

	go func() {
		items, err := r.fetchItems(ctx, false)
		nonTradableCh <- result{items: items, err: err}
	}()

	tradableRes := <-tradableCh
	nonTradableRes := <-nonTradableCh

	if tradableRes.err != nil {
		return nil, fmt.Errorf("fetch tradable items: %w", tradableRes.err)
	}
	if nonTradableRes.err != nil {
		return nil, fmt.Errorf("fetch non-tradable items: %w", nonTradableRes.err)
	}

	return r.mergeItems(tradableRes.items, nonTradableRes.items), nil
}

// fetchItems fetches items from Skinport API with the given tradable flag.
func (r *SkinportRepo) fetchItems(ctx context.Context, tradable bool) ([]skinportItem, error) {
	u, err := url.Parse(r.baseURL + "/items")
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	q := u.Query()
	q.Set("app_id", strconv.Itoa(r.appID))
	q.Set("currency", r.currency)
	if tradable {
		q.Set("tradable", "1")
	} else {
		q.Set("tradable", "0")
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept-Encoding", "br")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var reader io.Reader = resp.Body

	reader = brotli.NewReader(resp.Body)

	fmt.Println(resp.Body)

	var items []skinportItem
	if err := json.NewDecoder(reader).Decode(&items); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return items, nil
}

// mergeItems combines tradable and non-tradable items by market_hash_name.
func (r *SkinportRepo) mergeItems(tradable, nonTradable []skinportItem) []entity.Item {
	// Build map from non-tradable items
	nonTradableMap := make(map[string]*float64, len(nonTradable))
	for _, item := range nonTradable {
		nonTradableMap[item.MarketHashName] = item.MinPrice
	}

	// Build result using tradable items as base
	itemsMap := make(map[string]*entity.Item, len(tradable))

	for _, item := range tradable {
		itemsMap[item.MarketHashName] = &entity.Item{
			MarketHashName:      item.MarketHashName,
			Currency:            item.Currency,
			SuggestedPrice:      item.SuggestedPrice,
			ItemPage:            item.ItemPage,
			MarketPage:          item.MarketPage,
			MinPriceTradable:    item.MinPrice,
			MinPriceNonTradable: nonTradableMap[item.MarketHashName],
		}
	}

	// Add items that exist only in non-tradable
	for _, item := range nonTradable {
		if _, exists := itemsMap[item.MarketHashName]; !exists {
			itemsMap[item.MarketHashName] = &entity.Item{
				MarketHashName:      item.MarketHashName,
				Currency:            item.Currency,
				SuggestedPrice:      item.SuggestedPrice,
				ItemPage:            item.ItemPage,
				MarketPage:          item.MarketPage,
				MinPriceTradable:    nil,
				MinPriceNonTradable: item.MinPrice,
			}
		}
	}

	// Convert map to slice
	result := make([]entity.Item, 0, len(itemsMap))
	for _, item := range itemsMap {
		result = append(result, *item)
	}

	return result
}
