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

		mockCache.On("AllItemIds", mock.Anything).Return(cacheIds, nil)
		mockCache.On("GetItemsByIds", mock.Anything, cacheIds).Return(cachedItems)
		mockStorage.On("AllItemsByItemIds", mock.Anything, []int{3}).Return(storageItems, nil)
		mockCache.On("SetItemsByIds", mock.Anything, mock.Anything).Return()

		// When
		items, err := itemService.AllItems(context.TODO())

		// Then
		assert.NoError(t, err)
		assert.Len(t, items, 3)
		assert.Equal(t, 1, items[0].Id)
		assert.Equal(t, 2, items[1].Id)
		assert.Equal(t, 3, items[2].Id)

		mockCache.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	t.Run("cache doesn't have item IDs, fetch all items from storage", func(t *testing.T) {
		// Given
		storageItems := []storage.Item{
			{Id: 1, Name: "Item 1", Price: 100},
			{Id: 2, Name: "Item 2", Price: 200},
		}

		mockCache.On("AllItemIds", mock.Anything).Return(nil, cache.ErrorIdsNotFound)
		mockStorage.On("AllItems", mock.Anything).Return(storageItems, nil)
		mockCache.On("SetItemsByIds", mock.Anything, mock.Anything).Return()

		// When
		items, err := itemService.AllItems(context.TODO())

		// Then
		assert.NoError(t, err)
		assert.Len(t, items, 2)
		assert.Equal(t, 1, items[0].Id)
		assert.Equal(t, 2, items[1].Id)

		mockCache.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	t.Run("cache and storage both fail", func(t *testing.T) {
		// Given
		mockCache.On("AllItemIds", mock.Anything).Return(nil, cache.ErrorIdsNotFound)
		mockStorage.On("AllItems", mock.Anything).Return(nil, fmt.Errorf("storage error"))

		// When
		items, err := itemService.AllItems(context.TODO())

		// Then
		assert.Error(t, err)
		assert.Nil(t, items)

		mockCache.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})
}
