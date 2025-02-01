package storage_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mohammadne/ice-global/internal/entities"
	"github.com/mohammadne/ice-global/internal/storage"
)

func TestCarts(t *testing.T) {
	carts := storage.NewCarts(database)

	t.Run("retrieve-cart-by-cookie-and-status", func(t *testing.T) {
		ccokie := ""
		status := entities.CartStatusOpen
		result, err := carts.RetrieveCartByCookieAndStatus(context.TODO(), ccokie, status)
		if err != nil {
			t.Fatal(err, "error on retrieving cart")
		}

		object, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(object))
	})
}
