package services

import (
	"context"
	"fmt"

	"github.com/mohammadne/ice-global/internal/entities"
	"github.com/mohammadne/ice-global/internal/storage"
)

type Carts interface {
	AllCartItemsByUserId(ctx context.Context, userId int) ([]entities.CartItem, error)
}

func NewCarts(cartItemsStorage storage.CartItems, cartsStorage storage.Carts) Carts {
	return &carts{cartItemsStorage: cartItemsStorage, cartsStorage: cartsStorage}
}

type carts struct {
	cartItemsStorage storage.CartItems
	cartsStorage     storage.Carts
	itemsStorage     storage.Items
}

func (c *carts) AllCartItemsByUserId(ctx context.Context, userId int) ([]entities.CartItem, error) {
	storageCart, err := c.cartsStorage.RetrieveCartByUserIdAndStatus(ctx, userId, entities.CartStatusOpen)
	if err != nil {
		if err == storage.ErrorCartNotFound {
			return []entities.CartItem{}, nil
		}
		return nil, fmt.Errorf("error retrieving cart by user-id: %v", err)
	}

	storageCartItems, err := c.cartItemsStorage.AllCartItemsByCartId(ctx, storageCart.Id)
	if err != nil {
		if err == storage.ErrorCartItemNotFound {
			return []entities.CartItem{}, nil
		}
		return nil, fmt.Errorf("error retrieving cart-items by cart-id: %v", err)
	}

	itemIds := make([]int, 0, len(storageCartItems))
	for _, storageCartItem := range storageCartItems {
		itemIds = append(itemIds, storageCartItem.ItemId)
	}

	storageItems, err := c.itemsStorage.AllItemsByItemIds(ctx, itemIds)
	if err != nil {
		if err == storage.ErrorCartItemNotFound {
			return []entities.CartItem{}, nil
		}
		return nil, fmt.Errorf("error retrieving items by item-ids: %v", err)
	}

	cartItems := make([]entities.CartItem, 0, len(storageCartItems))
	for _, storageCartItem := range storageCartItems {
		cartItem := entities.CartItem{
			Id:        storageCartItem.Id,
			Cart:      &entities.Cart{},
			Quantity:  storageCartItem.Quantity,
			IsDeleted: storageCartItem.DeletedAt.Valid,
		}

		var item *entities.Item
		for _, storageItem := range storageItems {
			if storageCartItem.ItemId == storageItem.Id {
				item = &entities.Item{
					Id:    storageItem.Id,
					Name:  storageItem.Name,
					Price: storageItem.Price,
				}
			}
		}
		cartItem.Item = item

		cartItems = append(cartItems, cartItem)
	}

	return cartItems, nil
}
