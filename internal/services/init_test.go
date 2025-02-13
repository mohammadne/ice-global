package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/mohammadne/shopping-cart-manager/internal/entities"
	"github.com/mohammadne/shopping-cart-manager/internal/repositories/storage"
	"github.com/stretchr/testify/mock"
)

type MockItemsCache struct{ mock.Mock }

func (m *MockItemsCache) AllItemIds(ctx context.Context) ([]int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1) // Return nil and the error
	}
	return args.Get(0).([]int), args.Error(1)
}

func (m *MockItemsCache) GetItemsByIds(ctx context.Context, ids []int) map[int]entities.Item {
	args := m.Called(ctx, ids)
	return args.Get(0).(map[int]entities.Item)
}

func (m *MockItemsCache) SetItemsByIds(ctx context.Context, items []entities.Item) {
	m.Called(ctx, items)
}

type MockCartItemsStorage struct{ mock.Mock }

func (m *MockCartItemsStorage) AllCartItemsByCartId(ctx context.Context, cartId int) ([]storage.CartItem, error) {
	args := m.Called(ctx, cartId)
	if args.Get(0) == nil {
		return nil, args.Error(1) // Return nil and the error
	}
	return args.Get(0).([]storage.CartItem), args.Error(1)
}

func (m *MockCartItemsStorage) RetrieveCartItemByCartIdAndItemId(ctx context.Context, cartId, itemId int) (*storage.CartItem, error) {
	args := m.Called(ctx, cartId, itemId)
	if args.Get(0) == nil {
		return nil, args.Error(1) // Return nil and the error
	}
	return args.Get(0).(*storage.CartItem), args.Error(1)
}

func (m *MockCartItemsStorage) CreateCartItem(ctx context.Context, cartItem *storage.CartItem) (int, error) {
	args := m.Called(ctx, cartItem)
	if args.Get(0) == nil {
		return -1, args.Error(1) // Return nil and the error
	}
	return args.Int(0), args.Error(1)
}

func (m *MockCartItemsStorage) UpdateCartItem(ctx context.Context, cartItem *storage.CartItem) error {
	args := m.Called(ctx, cartItem)
	return args.Error(0)
}

func (m *MockCartItemsStorage) DeleteCartItemById(ctx context.Context, cartItemId int, time time.Time) error {
	args := m.Called(ctx, cartItemId, time)
	return args.Error(0)
}

type MockCartsStorage struct{ mock.Mock }

func (m *MockCartsStorage) RetrieveCartByCookieAndStatus(ctx context.Context, cookie string, status entities.CartStatus) (*storage.Cart, error) {
	args := m.Called(ctx, cookie, status)
	if args.Get(0) == nil {
		return nil, args.Error(1) // Return nil and the error
	}
	return args.Get(0).(*storage.Cart), args.Error(1)
}

func (m *MockCartsStorage) CreateCart(ctx context.Context, cart *storage.Cart) (int, error) {
	args := m.Called(ctx, cart)
	return args.Int(0), args.Error(1)
}

func (m *MockCartsStorage) RetrieveCartById(ctx context.Context, cartId int) (*storage.Cart, error) {
	args := m.Called(ctx, cartId)
	if args.Get(0) == nil {
		return nil, args.Error(1) // Return nil and the error
	}
	return args.Get(0).(*storage.Cart), args.Error(1)
}

type MockItemsStorage struct{ mock.Mock }

func (m *MockItemsStorage) AllItems(ctx context.Context) ([]storage.Item, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1) // Return nil and the error
	}
	return args.Get(0).([]storage.Item), args.Error(1)
}

func (m *MockItemsStorage) AllItemsByItemIds(ctx context.Context, ids []int) ([]storage.Item, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]storage.Item), args.Error(1)
}

func TestMain(m *testing.M) {
	m.Run()
}
