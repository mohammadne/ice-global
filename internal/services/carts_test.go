package services_test

import (
	"context"
	"testing"

	"github.com/mohammadne/ice-global/internal/entities"
	"github.com/mohammadne/ice-global/internal/repositories/storage"
	"github.com/mohammadne/ice-global/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCartsService(t *testing.T) {
	mockCartItemsStorage := new(MockCartItemsStorage)
	mockCartsStorage := new(MockCartsStorage)
	mockItemsStorage := new(MockItemsStorage)
	service := services.NewCarts(mockCartItemsStorage, mockCartsStorage, mockItemsStorage)

	t.Run("RetrieveCartOptional - cart found", func(t *testing.T) {

		mockCartsStorage.
			On("RetrieveCartByCookieAndStatus", mock.Anything, "cookie", entities.CartStatusOpen).
			Return(&storage.Cart{
				ID:     1,
				Cookie: "cookie",
			}, nil).Once()

		result, err := service.RetrieveCartOptional(context.TODO(), "cookie")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Id)
		mockCartsStorage.AssertExpectations(t)
	})

	t.Run("RetrieveCartOptional - cart not found", func(t *testing.T) {
		mockCartsStorage.On("RetrieveCartByCookieAndStatus", mock.Anything, "cookie", entities.CartStatusOpen).
			Return(nil, storage.ErrorCartNotFound).Once()

		result, err := service.RetrieveCartOptional(context.TODO(), "cookie")
		assert.NoError(t, err)
		assert.Nil(t, result)
		mockCartsStorage.AssertExpectations(t)
	})

	t.Run("RetrieveCartRequired - cart found", func(t *testing.T) {
		mockCartsStorage.On("RetrieveCartByCookieAndStatus", mock.Anything, "cookie", entities.CartStatusOpen).
			Return(&storage.Cart{
				ID:     1,
				Cookie: "cookie",
			}, nil).Once()

		result, err := service.RetrieveCartRequired(context.TODO(), "cookie")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Id)
		mockCartsStorage.AssertExpectations(t)
	})

	t.Run("RetrieveCartRequired - cart not found, create new cart", func(t *testing.T) {
		mockCartsStorage.On("RetrieveCartByCookieAndStatus", mock.Anything, "cookie", entities.CartStatusOpen).
			Return(nil, storage.ErrorCartNotFound).Once()
		mockCartsStorage.On("CreateCart", mock.Anything, mock.Anything).Return(1, nil).Once()

		result, err := service.RetrieveCartRequired(context.TODO(), "cookie")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Id)
		mockCartsStorage.AssertExpectations(t)
	})

	t.Run("AllCartItemsByCartId - items found", func(t *testing.T) {
		mockCartItemsStorage.On("AllCartItemsByCartId", mock.Anything, 1).Return([]storage.CartItem{
			{ID: 1, CartID: 1, Quantity: 2, ItemID: 1},
		}, nil).Once()
		mockItemsStorage.On("AllItemsByItemIds", mock.Anything, []int{1}).Return([]storage.Item{
			{ID: 1, Name: "Item 1", Price: 100},
		}, nil).Once()

		result, err := service.AllCartItemsByCartId(context.TODO(), 1)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		mockCartItemsStorage.AssertExpectations(t)
		mockItemsStorage.AssertExpectations(t)
	})

	t.Run("AddItemToCart - new item", func(t *testing.T) {
		mockCartItemsStorage.On("RetrieveCartItemByCartIdAndItemId", mock.Anything, 1, 1).
			Return(&storage.CartItem{}, storage.ErrorCartItemNotFound).Once()
		mockCartItemsStorage.On("CreateCartItem", mock.Anything, mock.Anything).Return(1, nil).Once()

		err := service.AddItemToCart(context.TODO(), 1, 1, 3)
		assert.NoError(t, err)
		mockCartItemsStorage.AssertExpectations(t)
	})

	t.Run("AddItemToCart - item exists", func(t *testing.T) {
		mockCartItemsStorage.On("RetrieveCartItemByCartIdAndItemId", mock.Anything, 1, 1).
			Return(&storage.CartItem{ID: 1, Quantity: 1}, nil).Once()
		mockCartItemsStorage.On("UpdateCartItem", mock.Anything, mock.Anything).Return(nil).Once()

		err := service.AddItemToCart(context.TODO(), 1, 1, 3)
		assert.NoError(t, err)
		mockCartItemsStorage.AssertExpectations(t)
	})

	t.Run("DeleteCartItem - cart closed", func(t *testing.T) {
		mockCartsStorage.On("RetrieveCartById", mock.Anything, 1).
			Return(&storage.Cart{Status: string(entities.CartStatusClosed)}, nil).Once()

		err := service.DeleteCartItem(context.TODO(), 1, 1)
		assert.Error(t, err)
		assert.Equal(t, services.ErrorCartHasBeenClosed, err)
		mockCartsStorage.AssertExpectations(t)
	})

	t.Run("DeleteCartItem - item deleted", func(t *testing.T) {
		mockCartsStorage.On("RetrieveCartById", mock.Anything, 1).
			Return(&storage.Cart{Status: string(entities.CartStatusOpen)}, nil).Once()
		mockCartItemsStorage.On("DeleteCartItemById", mock.Anything, 1, mock.Anything).Return(nil).Once()

		err := service.DeleteCartItem(context.TODO(), 1, 1)
		assert.NoError(t, err)
		mockCartsStorage.AssertExpectations(t)
		mockCartItemsStorage.AssertExpectations(t)
	})
}
