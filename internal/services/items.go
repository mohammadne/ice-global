package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mohammadne/shopping-cart-manager/internal/entities"
	"github.com/mohammadne/shopping-cart-manager/internal/repositories/cache"
	"github.com/mohammadne/shopping-cart-manager/internal/repositories/storage"
)

type Items interface {
	AllItems(ctx context.Context) ([]entities.Item, error)
}

func NewItems(itemsCache cache.Items, itemsStorage storage.Items) Items {
	return &items{itemsCache: itemsCache, itemsStorage: itemsStorage}
}

type items struct {
	itemsCache   cache.Items
	itemsStorage storage.Items
}

func (i *items) AllItems(ctx context.Context) (result []entities.Item, err error) {
	ids, err := i.itemsCache.AllItemIds(ctx)
	if err != nil {
		if err == cache.ErrorIdsNotFound {
			storageItems, err := i.itemsStorage.AllItems(ctx)
			if err != nil {
				return nil, fmt.Errorf("error retrieving items: %v", err)
			}

			result = make([]entities.Item, 0, len(storageItems))
			for _, storageItem := range storageItems {
				result = append(result, entities.Item{
					Id:    storageItem.ID,
					Name:  storageItem.Name,
					Price: storageItem.Price,
				})
			}

			i.itemsCache.SetItemsByIds(ctx, result)

			return result, nil
		}
		slog.Error("error retrieving items", "Err", err)
	}

	// retrieve from cached items
	resultMap := i.itemsCache.GetItemsByIds(ctx, ids)

	// here we retrieve only missed itemd from our mysql database
	// and then we store the retrieved items into the cache
	if count := len(ids) - len(resultMap); count > 0 {
		missedIds := make([]int, 0, count)
		for _, id := range ids {
			if _, exists := resultMap[id]; !exists {
				missedIds = append(missedIds, id)
			}
		}

		storageItems, err := i.itemsStorage.AllItemsByItemIds(ctx, missedIds)
		if err != nil {
			return nil, fmt.Errorf("error retrieving items by ids: %v", err)
		}

		result = make([]entities.Item, 0, len(ids))
		for _, storageItem := range storageItems {
			result = append(result, entities.Item{
				Id:    storageItem.ID,
				Name:  storageItem.Name,
				Price: storageItem.Price,
			})
		}

		i.itemsCache.SetItemsByIds(ctx, result)
	}

	// build the result, focus that in the pervious statements
	// we only retrieve missed items from cache so here we have
	// to append retrieved items from cache
	for _, item := range resultMap {
		result = append(result, item)
	}

	return result, nil
}
