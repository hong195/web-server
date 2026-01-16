package items

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hong195/web-server/internal/entity"
	"github.com/hong195/web-server/internal/repo"
	"github.com/hong195/web-server/pkg/cache"
	"github.com/hong195/web-server/pkg/logger"
)

const cacheKey = "skinport:items"

// UseCase implements usecase.Items interface.
type UseCase struct {
	repo   repo.ItemsRepo
	cache  cache.Cache
	logger logger.Interface
	ttl    time.Duration
}

// New creates a new Items usecase.
func New(repo repo.ItemsRepo, cache cache.Cache, logger logger.Interface, ttlSec int) *UseCase {
	return &UseCase{
		repo:   repo,
		cache:  cache,
		logger: logger,
		ttl:    time.Duration(ttlSec) * time.Second,
	}
}

// StartBackgroundRefresh starts background cache refresh.
// It immediately loads data and then refreshes every ttl interval.
func (uc *UseCase) StartBackgroundRefresh(ctx context.Context) {
	// Initial load
	uc.refresh(ctx)

	// Periodic refresh, кэш всегда на готове, повышает доступность
	ticker := time.NewTicker(uc.ttl)
	go func() {
		for {
			select {
			case <-ticker.C:
				uc.refresh(ctx)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

// refresh fetches items from repo and updates cache.
func (uc *UseCase) refresh(ctx context.Context) {
	items, err := uc.repo.GetItems(ctx)
	if err != nil {
		uc.logger.Error("failed to refresh items cache: %v", err)
		return
	}

	data, err := json.Marshal(items)
	if err != nil {
		uc.logger.Error("failed to marshal items: %v", err)
		return
	}

	uc.cache.Set(cacheKey, data, uc.ttl)
	uc.logger.Info("items cache refreshed, count: %d", len(items))
}

// GetItems returns items from cache or fetches from repo.
func (uc *UseCase) GetItems(ctx context.Context) ([]entity.Item, error) {
	// Try cache first
	if cached, ok := uc.cache.Get(cacheKey); ok {
		var items []entity.Item
		if err := json.Unmarshal(cached, &items); err == nil {
			return items, nil
		}
	}

	// Fallback, На случай если в кэше нет, но таких ситуаций не должно быть
	items, err := uc.repo.GetItems(ctx)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if data, err := json.Marshal(items); err == nil {
		uc.cache.Set(cacheKey, data, uc.ttl)
	}

	return items, nil
}
