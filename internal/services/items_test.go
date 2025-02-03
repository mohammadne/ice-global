package services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/mohammadne/ice-global/internal/entities"
	"github.com/mohammadne/ice-global/internal/repositories/cache"
	"github.com/mohammadne/ice-global/internal/repositories/storage"
	"github.com/mohammadne/ice-global/internal/services"
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
			{Id: 1, Name: "Item 1", Price: 100},
			{Id: 2, Name: "Item 2", Price: 200},
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
			{Id: 3, Name: "Item 3", Price: 300},
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
