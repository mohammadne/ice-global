package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mohammadne/ice-global/internal/entities"
	"github.com/mohammadne/ice-global/internal/storage"
)

type Carts interface {
	RetrieveCartOptional(ctx context.Context, cookie string) (*entities.Cart, error)
	RetrieveCartRequired(ctx context.Context, cookie string) (*entities.Cart, error)
	AllCartItemsByCartId(ctx context.Context, cartId int) ([]entities.CartItem, error)
}

func NewCarts(cartItemsStorage storage.CartItems, cartsStorage storage.Carts) Carts {
	return &carts{cartItemsStorage: cartItemsStorage, cartsStorage: cartsStorage}
}

type carts struct {
	cartItemsStorage storage.CartItems
	cartsStorage     storage.Carts
	itemsStorage     storage.Items
}

func (u *carts) RetrieveCartOptional(ctx context.Context, cookie string) (result *entities.Cart, err error) {
	{ // validation
		if cookie == "" {
			return nil, errors.New("the cookie should be provided")
		}
	}

	storageUser, err := u.cartsStorage.RetrieveCartByCookieAndStatus(ctx, cookie, entities.CartStatusOpen)
	if err != nil {
		if err == storage.ErrorCartNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error retrieving cart: %v", err)
	}

	return &entities.Cart{
		Id:     storageUser.Id,
		Cookie: storageUser.Cookie,
	}, nil
}

func (u *carts) RetrieveCartRequired(ctx context.Context, cookie string) (result *entities.Cart, err error) {
	{ // validation
		if cookie == "" {
			return nil, errors.New("the cookie should be provided")
		}
	}

	storageCart, err := u.cartsStorage.RetrieveCartByCookieAndStatus(ctx, cookie, entities.CartStatusOpen)
	if err != nil {
		if err == storage.ErrorCartNotFound {
			storageCart := &storage.Cart{
				Cookie:    cookie,
				Status:    entities.CartStatusOpen,
				CreatedAt: time.Now(),
			}

			id, err := u.cartsStorage.CreateCart(ctx, storageCart)
			if err != nil {
				return nil, fmt.Errorf("error creating cart: %v", err)
			}
			return &entities.Cart{
				Id:     id,
				Cookie: cookie,
				Status: entities.CartStatusOpen,
			}, nil
		}
		return nil, fmt.Errorf("error retrieving cart: %v", err)
	}

	return &entities.Cart{
		Id:     storageCart.Id,
		Cookie: storageCart.Cookie,
		Status: entities.CartStatus(storageCart.Status),
	}, nil
}

func (c *carts) AllCartItemsByCartId(ctx context.Context, cartId int) ([]entities.CartItem, error) {
	// storageCart, err := c.cartsStorage.RetrieveCartByCookieAndStatus(ctx, cookie, entities.CartStatusOpen)
	// if err != nil {
	// 	if err == storage.ErrorCartNotFound {
	// 		return []entities.CartItem{}, nil
	// 	}
	// 	return nil, fmt.Errorf("error retrieving cart by user-id: %v", err)
	// }

	storageCartItems, err := c.cartItemsStorage.AllCartItemsByCartId(ctx, cartId)
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
