package functional

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	// "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/mohammadne/ice-global/internal/entities"
	"github.com/mohammadne/ice-global/internal/repositories/cache"
	"github.com/mohammadne/ice-global/pkg/redis"
)

var itemCache cache.Items

func setupRedis() {
	cfg := redis.Config{Address: "localhost:4000", DB: 1}
	redisClient, err := redis.Open(&cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open redis: %v\n", err)
		os.Exit(1) // Exit with a non-zero status code
	}

	itemCache = cache.NewItems(redisClient)
}

func TestItemsSetGet(t *testing.T) {
	setupRedis() // Ensure Redis client is set up

	item := entities.Item{
		Id:    1,
		Name:  "Test Item",
		Price: 100,
	}

	itemCache.SetItemsByIds(context.Background(), []entities.Item{item})

	{ // test AllItemIds
		allItemIds, err := itemCache.AllItemIds(context.Background())
		if err != nil && !errors.Is(err, cache.ErrorIdsNotFound) {
			t.Fatalf("an error occured during retrieving all item ids, %v", err)
		}

		// Assert that the cached item is the same as the added item
		assert.Len(t, allItemIds, 1)
		assert.Equal(t, allItemIds[0], item.Id)
	}

	{ // test GetItemsByIds
		cachedItems := itemCache.GetItemsByIds(context.Background(), []int{1})

		// Assert that the cached item is the same as the added item
		assert.Len(t, cachedItems, 1)
		assert.Equal(t, cachedItems[1].Id, item.Id)
		assert.Equal(t, cachedItems[1].Name, item.Name)
		assert.Equal(t, cachedItems[1].Price, item.Price)
	}
}
