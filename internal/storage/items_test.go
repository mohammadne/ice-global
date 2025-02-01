package storage_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mohammadne/ice-global/internal/storage"
)

func TestItems(t *testing.T) {
	carts := storage.NewItems(mysqlDatabase, redisDatabase)

	t.Run("retrieve-cart-by-cookie-and-status", func(t *testing.T) {
		result, err := carts.AllItems(context.TODO())
		if err != nil {
			t.Fatal(err, "error on retrieving cart")
		}

		object, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(object))
	})
}
