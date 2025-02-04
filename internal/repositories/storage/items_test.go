package storage_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var itemColumns = []string{
	"id",
	"name",
	"price",
	"created_at",
	"updated_at",
}

func TestAllItems(t *testing.T) {
	mockDB.
		ExpectQuery("SELECT id, name, price, created_at, updated_at FROM items").
		WillReturnRows(sqlmock.NewRows(itemColumns).
			AddRow(1, "shoe", 100, time.Now(), nil).
			AddRow(2, "basket", 200, time.Now(), nil),
		)

	items, err := itemsStorage.AllItems(context.TODO())
	if err != nil {
		t.Errorf("expect no errors %v", err)
	}

	if len(items) != 2 {
		t.Error("invalid items length result")
	} else if items[0].ID != 1 {
		t.Error("invalid item-id has been returned")
	}

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestAllItemsByItemIds(t *testing.T) {
	itemIds := []int{1, 2}

	t.Run("with_empty_result", func(t *testing.T) {
		mockDB.
			ExpectQuery(regexp.QuoteMeta(("SELECT id, name, price, created_at, updated_at FROM items WHERE id IN (?, ?)"))).
			WithArgs(1, 2).
			WillReturnRows(sqlmock.NewRows(itemColumns))

		items, err := itemsStorage.AllItemsByItemIds(context.TODO(), itemIds)
		if err != nil {
			t.Errorf("expect no errors %v", err)
		}

		if len(items) != 0 {
			t.Error("expecting no result")
		}

		if err := mockDB.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("with_full_result", func(t *testing.T) {
		mockDB.
			ExpectQuery(regexp.QuoteMeta(("SELECT id, name, price, created_at, updated_at FROM items WHERE id IN (?, ?)"))).
			WithArgs(1, 2).
			WillReturnRows(sqlmock.NewRows(itemColumns).
				AddRow(1, "shoe", 100, time.Now(), nil).
				AddRow(2, "basket", 200, time.Now(), nil),
			)

		items, err := itemsStorage.AllItemsByItemIds(context.TODO(), itemIds)
		if err != nil {
			t.Errorf("expect no errors %v", err)
		}

		if len(items) != len(itemIds) {
			t.Error("invalid items length result")
		} else if items[0].ID != 1 {
			t.Error("invalid item-id has been returned")
		}

		if err := mockDB.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})
}
