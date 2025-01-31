package services

import (
	"context"
	"fmt"

	"github.com/mohammadne/ice-global/internal/entities"
	"github.com/mohammadne/ice-global/internal/storage"
)

type Items interface {
	AllItems(ctx context.Context) ([]entities.Item, error)
}

func NewItems(itemsStorage storage.Items) Items {
	return &items{itemsStorage: itemsStorage}
}

type items struct {
	itemsStorage storage.Items
}

func (i *items) AllItems(ctx context.Context) (result []entities.Item, err error) {
	storageItems, err := i.itemsStorage.AllItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving items: %v", err)
	}

	result = make([]entities.Item, 0, len(storageItems))
	for _, storageItem := range storageItems {
		result = append(result, entities.Item{
			Id:    storageItem.Id,
			Name:  storageItem.Name,
			Price: storageItem.Price,
		})
	}

	return
}
