package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/mohammadne/ice-global/internal/entities"
	"github.com/mohammadne/ice-global/pkg/redis"
)

type Items interface {
	// AllItemIds retrieves all the item IDs
	AllItemIds(ctx context.Context) ([]int, error)

	// GetItemsByIds
	GetItemsByIds(ctx context.Context, ids []int) map[int]entities.Item

	// SetItemsByIds
	SetItemsByIds(ctx context.Context, items []entities.Item)
}

func NewItems(redis *redis.Redis) Items {
	return &items{redis: redis}
}

type items struct {
	redis *redis.Redis
}

var (
	ErrorIdsNotFound = errors.New("no ids have been found")
	idsKey           = "item:all:ids" // Key that stores all item IDs
)

// AllItemIds retrieves all item IDs from Redis.
func (i *items) AllItemIds(ctx context.Context) ([]int, error) {
	ids, err := i.redis.SMembers(ctx, idsKey).Result()
	if errors.Is(err, redis.Nil) || len(ids) == 0 {
		return []int{}, ErrorIdsNotFound
	} else if err != nil {
		return nil, fmt.Errorf("error retrieving item IDs from Redis: %v", err)
	}

	itemIds := make([]int, 0, len(ids))
	for _, idString := range ids {
		id, err := strconv.Atoi(idString)
		if err != nil {
			slog.Error("error converting string to int for item Id", "Err", err)
			continue
		}
		itemIds = append(itemIds, id)
	}

	return itemIds, nil
}

var itemKeyPrefix = "item:%d" // Key that stores item values

// GetItemsByIds retrieves cached items from Redis by their IDs.
func (c *items) GetItemsByIds(ctx context.Context, ids []int) map[int]entities.Item {
	result := make(map[int]entities.Item, len(ids))

	for _, id := range ids {
		cacheKey := fmt.Sprintf(itemKeyPrefix, id)

		// Retrieve item from Redis by ID
		cachedItem, err := c.redis.Get(ctx, cacheKey).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			slog.Error("error checking Redis for item", slog.Int("id", id), "Err", err)
			continue
		}

		// Cache hit: Unmarshal the cached item
		var item entities.Item
		err = json.Unmarshal([]byte(cachedItem), &item)
		if err != nil {
			slog.Error("error unmarshalling cached item", slog.Int("id", id), "Err", err)
			continue
		}
		result[id] = item
	}

	return result
}

var cacheTTL = 30 * time.Minute // Time to live for cached items

// SetItemsByIds caches the given items in Redis by their ID.
func (c *items) SetItemsByIds(ctx context.Context, items []entities.Item) {
	for _, item := range items {
		marshaledItem, err := json.Marshal(item)
		if err != nil {
			slog.Error("error marshalling item", slog.Int("id", item.Id), "Err", err)
			continue
		}

		// Set the item in Redis with TTL
		cacheKey := fmt.Sprintf(itemKeyPrefix, item.Id)
		err = c.redis.Set(ctx, cacheKey, marshaledItem, cacheTTL).Err()
		if err != nil {
			slog.Error("error caching item", slog.Int("id", item.Id), "Err", err)
			continue
		}

		// Add item Id to the set of all item Ids
		err = c.redis.SAdd(ctx, idsKey, item.Id).Err()
		if err != nil {
			slog.Error("error adding item Id to set", slog.Int("id", item.Id), "Err", err)
			continue
		}
	}
}
