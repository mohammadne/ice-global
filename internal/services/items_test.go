package services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/mohammadne/shopping-cart-manager/internal/entities"
	"github.com/mohammadne/shopping-cart-manager/internal/repositories/cache"
	"github.com/mohammadne/shopping-cart-manager/internal/repositories/storage"
	"github.com/mohammadne/shopping-cart-manager/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAllItems(t *testing.T) {
	mockCache := new(MockItemsCache)
	mockStorage := new(MockItemsStorage)
	itemService := services.NewItems(mockCache, mockStorage)

	t.Run("cache doesn't have item IDs, fetch all items from storage", func(t *testing.T) {

		// Given
		storageItems := []storage.Item{
			{ID: 1, Name: "Item 1", Price: 100},
			{ID: 2, Name: "Item 2", Price: 200},
		}

		mockCache.On("AllItemIds", mock.Anything).Return(nil, cache.ErrorIdsNotFound).Once()
		mockStorage.On("AllItems", mock.Anything).Return(storageItems, nil).Once()
		mockCache.On("SetItemsByIds", mock.Anything, mock.Anything).Return().Once()

		// When
		items, err := itemService.AllItems(context.TODO())

		// Then
		assert.NoError(t, err)
		assert.Len(t, items, 2)
	})

	t.Run("cache and storage both fail", func(t *testing.T) {
		// Given
		mockCache.On("AllItemIds", mock.Anything).Return(nil, cache.ErrorIdsNotFound).Once()
		mockStorage.On("AllItems", mock.Anything).Return(nil, fmt.Errorf("storage error")).Once()

		// When
		items, err := itemService.AllItems(context.TODO())

		// Then
		assert.Error(t, err)
		assert.Empty(t, items)
	})

	t.Run("cache has item IDs, fetch missing items from storage", func(t *testing.T) {
		// Given
		cacheIds := []int{1, 2, 3}
		cachedItems := map[int]entities.Item{
			1: {Id: 1, Name: "Item 1", Price: 100},
			2: {Id: 2, Name: "Item 2", Price: 200},
		}
		storageItems := []storage.Item{
			{ID: 3, Name: "Item 3", Price: 300},
		}

		mockCache.On("AllItemIds", mock.Anything).Return(cacheIds, nil).Once()
		mockCache.On("GetItemsByIds", mock.Anything, cacheIds).Return(cachedItems).Once()
		mockStorage.On("AllItemsByItemIds", mock.Anything, []int{3}).Return(storageItems, nil).Once()
		mockCache.On("SetItemsByIds", mock.Anything, mock.Anything).Return().Once()

		// When
		items, err := itemService.AllItems(context.TODO())

		// Then
		assert.NoError(t, err)
		assert.Len(t, items, 3)
	})
}

func BenchmarkAllItems(b *testing.B) {
	ctx := context.Background()

	mockCache := new(MockItemsCache)
	mockStorage := new(MockItemsStorage)

	// Mock behavior for cache miss
	mockCache.On("AllItemIds", ctx).Return([]int{}, cache.ErrorIdsNotFound)
	mockStorage.On("AllItems", ctx).Return([]storage.Item{
		{ID: 1, Name: "Item1", Price: 100},
		{ID: 2, Name: "Item2", Price: 200},
	}, nil)
	mockCache.On("SetItemsByIds", ctx, mock.Anything).Return()

	service := services.NewItems(mockCache, mockStorage)

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		_, _ = service.AllItems(ctx)
	}

	// Verify expectations
	mockCache.AssertExpectations(b)
	mockStorage.AssertExpectations(b)
}
