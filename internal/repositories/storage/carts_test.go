package storage_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mohammadne/ice-global/internal/entities"
	"github.com/mohammadne/ice-global/internal/repositories/storage"
)

var cartColumns = []string{
	"id",
	"session_id",
	"status",
	"created_at",
	"deleted_at",
}

func TestCreateCart(t *testing.T) {
	storageCart := storage.Cart{
		Cookie:    "sample-cookie-timestamp",
		Status:    "open",
		CreatedAt: time.Now(),
	}

	mockDB.
		ExpectExec(regexp.QuoteMeta("INSERT INTO cart_entities (session_id, status, created_at) VALUES (?, ?, ?)")).
		WithArgs(storageCart.Cookie, storageCart.Status, storageCart.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := cartsStorage.CreateCart(context.TODO(), &storageCart)
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

func TestRetrieveCartById(t *testing.T) {
	itemId := 1

	t.Run("with_empty_result", func(t *testing.T) {
		mockDB.
			ExpectQuery(regexp.QuoteMeta(("SELECT id, session_id, status, created_at, deleted_at FROM cart_entities	WHERE id = ?"))).
			WithArgs(itemId).
			WillReturnRows(sqlmock.NewRows(cartColumns))

		_, err := cartsStorage.RetrieveCartById(context.TODO(), itemId)
		if !errors.Is(err, storage.ErrorCartNotFound) {
			t.Errorf("expect ErrorCartNotFound error %v", err)
		}

		if err := mockDB.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("with_result", func(t *testing.T) {
		storageCart := storage.Cart{
			Id:        12,
			Cookie:    "sample-cookie-timestamp",
			Status:    "open",
			CreatedAt: time.Now(),
		}

		mockDB.
			ExpectQuery(regexp.QuoteMeta(("SELECT id, session_id, status, created_at, deleted_at FROM cart_entities	WHERE id = ?"))).
			WithArgs(itemId).
			WillReturnRows(sqlmock.NewRows(cartColumns).
				AddRow(storageCart.Id, storageCart.Cookie, storageCart.Status, storageCart.CreatedAt, nil))

		cart, err := cartsStorage.RetrieveCartById(context.TODO(), itemId)
		if err != nil {
			t.Errorf("expect no errors %v", err)
		}

		if cart == nil {
			t.Error("expecting result but cart is empty")
		} else if cart.Id != storageCart.Id {
			t.Error("invalid cart-id has been returned")
		}

		if err := mockDB.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})
}

func TestRetrieveCartByCookieAndStatus(t *testing.T) {
	cookie := "sample-cookie-timestamp"
	status := entities.CartStatusOpen

	t.Run("with_empty_result", func(t *testing.T) {
		mockDB.
			ExpectQuery(regexp.QuoteMeta(("SELECT id, session_id, status, created_at, deleted_at FROM cart_entities WHERE session_id = ? AND status = ?"))).
			WithArgs(cookie, status).
			WillReturnRows(sqlmock.NewRows(cartColumns))

		_, err := cartsStorage.RetrieveCartByCookieAndStatus(context.TODO(), cookie, status)
		if !errors.Is(err, storage.ErrorCartNotFound) {
			t.Errorf("expect ErrorCartNotFound error %v", err)
		}

		if err := mockDB.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("with_result", func(t *testing.T) {
		storageCart := storage.Cart{
			Id:        12,
			Cookie:    cookie,
			Status:    string(status),
			CreatedAt: time.Now(),
		}

		mockDB.
			ExpectQuery(regexp.QuoteMeta(("SELECT id, session_id, status, created_at, deleted_at FROM cart_entities WHERE session_id = ? AND status = ?"))).
			WithArgs(cookie, status).
			WillReturnRows(sqlmock.NewRows(cartColumns).
				AddRow(storageCart.Id, storageCart.Cookie, storageCart.Status, storageCart.CreatedAt, nil))

		cart, err := cartsStorage.RetrieveCartByCookieAndStatus(context.TODO(), cookie, status)
		if err != nil {
			t.Errorf("expect no errors %v", err)
		}

		if cart == nil {
			t.Error("expecting result but cart is empty")
		} else if cart.Id != storageCart.Id {
			t.Error("invalid cart-id has been returned")
		}

		if err := mockDB.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})
}
