package items

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/hong195/web-server/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRepo is a mock implementation of repo.ItemsRepo.
type mockRepo struct {
	items []entity.Item
	err   error
}

func (m *mockRepo) GetItems(_ context.Context) ([]entity.Item, error) {
	return m.items, m.err
}

// mockCache is a mock implementation of cache.Cache.
type mockCache struct {
	data map[string][]byte
}

func newMockCache() *mockCache {
	return &mockCache{data: make(map[string][]byte)}
}

func (m *mockCache) Get(key string) ([]byte, bool) {
	v, ok := m.data[key]
	return v, ok
}

func (m *mockCache) Set(key string, value []byte, _ time.Duration) {
	m.data[key] = value
}

// mockLogger is a mock implementation of logger.Interface.
type mockLogger struct{}

func (m *mockLogger) Debug(_ interface{}, _ ...interface{}) {}
func (m *mockLogger) Info(_ string, _ ...interface{})       {}
func (m *mockLogger) Warn(_ string, _ ...interface{})       {}
func (m *mockLogger) Error(_ interface{}, _ ...interface{}) {}
func (m *mockLogger) Fatal(_ interface{}, _ ...interface{}) {}

func TestGetItems(t *testing.T) {
	t.Parallel()

	tradablePrice := 10.5
	nonTradablePrice := 8.0
	errRepo := errors.New("connection refused")

	tests := []struct {
		name           string
		cacheData      []byte
		repoItems      []entity.Item
		repoErr        error
		wantItems      []entity.Item
		wantErr        error
		wantCacheAfter bool
	}{
		{
			name: "from cache",
			cacheData: mustMarshal([]entity.Item{
				{
					MarketHashName:      "AK-47 | Redline",
					MinPriceTradable:    &tradablePrice,
					MinPriceNonTradable: &nonTradablePrice,
				},
			}),
			repoItems: nil,
			repoErr:   errors.New("should not be called"),
			wantItems: []entity.Item{
				{
					MarketHashName:      "AK-47 | Redline",
					MinPriceTradable:    &tradablePrice,
					MinPriceNonTradable: &nonTradablePrice,
				},
			},
			wantErr:        nil,
			wantCacheAfter: true,
		},
		{
			name:      "cache miss fallback to repo",
			cacheData: nil,
			repoItems: []entity.Item{
				{
					MarketHashName:   "AWP | Asiimov",
					MinPriceTradable: &tradablePrice,
				},
			},
			repoErr: nil,
			wantItems: []entity.Item{
				{
					MarketHashName:   "AWP | Asiimov",
					MinPriceTradable: &tradablePrice,
				},
			},
			wantErr:        nil,
			wantCacheAfter: true,
		},
		{
			name:      "invalid cache data fallback to repo",
			cacheData: []byte("invalid json"),
			repoItems: []entity.Item{
				{
					MarketHashName:   "M4A4 | Howl",
					MinPriceTradable: &tradablePrice,
				},
			},
			repoErr: nil,
			wantItems: []entity.Item{
				{
					MarketHashName:   "M4A4 | Howl",
					MinPriceTradable: &tradablePrice,
				},
			},
			wantErr:        nil,
			wantCacheAfter: true,
		},
		{
			name:           "repo error",
			cacheData:      nil,
			repoItems:      nil,
			repoErr:        errRepo,
			wantItems:      nil,
			wantErr:        errRepo,
			wantCacheAfter: false,
		},
		{
			name:           "empty items",
			cacheData:      nil,
			repoItems:      []entity.Item{},
			repoErr:        nil,
			wantItems:      []entity.Item{},
			wantErr:        nil,
			wantCacheAfter: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cache := newMockCache()
			if tt.cacheData != nil {
				cache.Set(cacheKey, tt.cacheData, time.Minute)
			}

			repo := &mockRepo{items: tt.repoItems, err: tt.repoErr}
			uc := New(repo, cache, &mockLogger{}, 300)

			items, err := uc.GetItems(context.Background())

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, items)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantItems, items)
			}

			if tt.wantCacheAfter {
				_, ok := cache.Get(cacheKey)
				assert.True(t, ok)
			}
		})
	}
}

func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
