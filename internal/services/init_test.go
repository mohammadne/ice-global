package services_test

import (
	"context"
	"testing"

	"github.com/mohammadne/ice-global/internal/entities"
	"github.com/mohammadne/ice-global/internal/repositories/storage"
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
