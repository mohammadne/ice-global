package storage_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mohammadne/ice-global/internal/repositories/storage"
)

var cartItemColumns = []string{
	"id",
	"cart_id",
	"item_id",
	"quantity",
	"created_at",
	"updated_at",
	"deleted_at",
}

func TestCreateCartItem(t *testing.T) {
	storageCartItem := storage.CartItem{
		CartID:    1,
		ItemID:    1,
		Quantity:  2,
		CreatedAt: time.Now(),
	}

	mockDB.
		ExpectExec(regexp.QuoteMeta("INSERT INTO cart_items (cart_id, item_id, quantity, created_at) VALUES (?, ?, ?, ?)")).
		WithArgs(storageCartItem.CartID, storageCartItem.ItemID, storageCartItem.Quantity, storageCartItem.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := cartItemStorage.CreateCartItem(context.TODO(), &storageCartItem)
	if err != nil {
		t.Errorf("expect no errors %v", err)
	}

	if id != 1 {
		t.Error("invalid inserted cart-id has been returned")
	}

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestAllCartItemsByCartId(t *testing.T) {
	cartItemIds := []int{1, 2}

	t.Run("with_empty_result", func(t *testing.T) {
		mockDB.
			ExpectQuery("SELECT id, cart_id, item_id, quantity, created_at, updated_at, deleted_at FROM cart_items WHERE cart_id = ?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(cartItemColumns))

		cartItems, err := cartItemStorage.AllCartItemsByCartId(context.TODO(), 1)
		if err != nil {
			t.Errorf("expect no errors %v", err)
		}

		if len(cartItems) != 0 {
			t.Error("expecting no result")
		}

		if err := mockDB.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("with_full_result", func(t *testing.T) {
		mockDB.
			ExpectQuery("SELECT id, cart_id, item_id, quantity, created_at, updated_at, deleted_at FROM cart_items WHERE cart_id = ?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(cartItemColumns).
				AddRow(1, 1, 1, 1, time.Now(), nil, nil).
				AddRow(2, 1, 2, 3, time.Now(), nil, nil),
			)

		cartItems, err := cartItemStorage.AllCartItemsByCartId(context.TODO(), 1)
		if err != nil {
			t.Errorf("expect no errors %v", err)
		}

		if len(cartItems) != len(cartItemIds) {
			t.Error("invalid items length result")
		} else if cartItems[0].ID != 1 {
			t.Error("invalid item-id has been returned")
		}

		if err := mockDB.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})
}

func TestRetrieveCartItemByCartIdAndItemId(t *testing.T) {
	cartId := 1
	itemId := 1

	t.Run("with_empty_result", func(t *testing.T) {
		mockDB.
			ExpectQuery(regexp.QuoteMeta("SELECT id, cart_id, item_id, quantity, created_at, updated_at, deleted_at FROM cart_items WHERE cart_id = ? AND item_id = ?")).
			WithArgs(cartId, itemId).
			WillReturnRows(sqlmock.NewRows(cartItemColumns))

		_, err := cartItemStorage.RetrieveCartItemByCartIdAndItemId(context.TODO(), cartId, itemId)
		if !errors.Is(err, storage.ErrorCartItemNotFound) {
			t.Errorf("expect ErrorCartItemNotFound error: %v", err)
		}

		if err := mockDB.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("with_result", func(t *testing.T) {
		storageCartItem := storage.CartItem{
			ID:        12,
			CartID:    cartId,
			ItemID:    itemId,
			Quantity:  3,
			CreatedAt: time.Now(),
		}

		mockDB.
			ExpectQuery(regexp.QuoteMeta("SELECT id, cart_id, item_id, quantity, created_at, updated_at, deleted_at FROM cart_items WHERE cart_id = ? AND item_id = ?")).
			WithArgs(cartId, itemId).
			WillReturnRows(sqlmock.NewRows(cartItemColumns).
				AddRow(storageCartItem.ID, storageCartItem.CartID, storageCartItem.ItemID, storageCartItem.Quantity, storageCartItem.CreatedAt, nil, nil))

		cartItem, err := cartItemStorage.RetrieveCartItemByCartIdAndItemId(context.TODO(), cartId, itemId)
		if err != nil {
			t.Errorf("expect no errors %v", err)
		}

		if cartItem == nil {
			t.Error("expecting result but cart is empty")
		} else if cartItem.ID != storageCartItem.ID {
			t.Error("invalid cart-id has been returned")
		}

		if err := mockDB.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})
}

func TestUpdateCartItem(t *testing.T) {
	storageCartItem := storage.CartItem{
		ID:        1,
		CartID:    1,
		ItemID:    1,
		Quantity:  2,
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	mockDB.
		ExpectExec(regexp.QuoteMeta("UPDATE cart_items SET cart_id = ?, item_id = ?, quantity = ?, updated_at = ? WHERE id = ?")).
		WithArgs(storageCartItem.CartID, storageCartItem.ItemID, storageCartItem.Quantity, storageCartItem.UpdatedAt, storageCartItem.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := cartItemStorage.UpdateCartItem(context.TODO(), &storageCartItem)
	if err != nil {
		t.Errorf("expect no errors %v", err)
	}

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestDeleteCartItemById(t *testing.T) {
	deletedAt := time.Now()
	cartItemId := 1

	mockDB.
		ExpectExec(regexp.QuoteMeta("UPDATE cart_items SET deleted_at = ? WHERE id = ?")).
		WithArgs(deletedAt, cartItemId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := cartItemStorage.DeleteCartItemById(context.TODO(), cartItemId, deletedAt)
	if err != nil {
		t.Errorf("expect no errors %v", err)
	}

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
