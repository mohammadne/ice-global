package cache_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"

	"github.com/mohammadne/ice-global/internal/entities"
	"github.com/mohammadne/ice-global/internal/repositories/cache"
	"github.com/mohammadne/ice-global/pkg/redis"
)

func TestItems(t *testing.T) {
	s := miniredis.RunT(t)

	cfg := redis.Config{Address: s.Addr(), Timeout: time.Second * 2}
	c, err := redis.Open(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	itemsCache := cache.NewItems(c)

	testItems := []entities.Item{
		{
			Id:    1,
			Name:  "shoe",
			Price: 100,
		},
		{
			Id:    2,
			Name:  "shoe",
			Price: 200,
		},
		{
			Id:    3,
			Name:  "shoe",
			Price: 300,
		},
		{
			Id:    4,
			Name:  "shoe",
			Price: 400,
		},
	}

	t.Run("all_item_ids", func(t *testing.T) {
		var idsKey = "item:all:ids"

		t.Run("empty_keys", func(t *testing.T) {
			s.FlushAll()

			_, err := itemsCache.AllItemIds(context.TODO())
			if err != cache.ErrorIdsNotFound {
				t.Error("when no key exists, it should return ErrorIdsNotFound error")
			}
		})

		t.Run("check_length", func(t *testing.T) {
			s.FlushAll()

			idsToPass := []string{"1", "2"}
			s.SAdd(idsKey, idsToPass...)

			idsToGet, err := s.SMembers(idsKey)
			if err != nil {
				t.Error("should not return any error")
			}

			if len(idsToPass) != len(idsToGet) {
				t.Errorf("expect %d, got %d for ids length", len(idsToPass), len(idsToGet))
			}
		})
	})

	t.Run("get_items_by_ids", func(t *testing.T) {
		var (
			itemKeyPrefix = "item:%d"
			cacheTTL      = 3 * time.Second
		)

		t.Run("pass_no_ids", func(t *testing.T) {
			items := itemsCache.GetItemsByIds(context.TODO(), nil)
			if len(items) != 0 {
				t.Error("when passing no ids, no item should be returned")
			}
		})

		t.Run("get_items_not_set", func(t *testing.T) {
			s.FlushAll()

			items := itemsCache.GetItemsByIds(context.TODO(), []int{1, 2, 3})
			if len(items) != 0 {
				t.Error("when setting no items, the result should be empty")
			}
		})

		t.Run("set_items_with_ttl", func(t *testing.T) {
			s.FlushAll()

			for index := 1; index < 3; index++ {
				marshaledItem, err := json.Marshal(testItems[index])
				if err != nil {
					t.Errorf("error marshalling item %v", err)
					continue
				}

				key := fmt.Sprintf(itemKeyPrefix, testItems[index].Id)
				s.Set(key, string(marshaledItem))
				s.SetTTL(key, cacheTTL)
			}

			items := itemsCache.GetItemsByIds(context.TODO(), []int{1, 2, 3})
			if len(items) != 2 {
				t.Error("only item id 2 and 3 should be returned")
			}
		})

		t.Run("check_ttl", func(t *testing.T) {
			s.FlushAll()

			for index := 1; index < 3; index++ {
				marshaledItem, err := json.Marshal(testItems[index])
				if err != nil {
					t.Errorf("error marshalling item %v", err)
					continue
				}

				key := fmt.Sprintf(itemKeyPrefix, testItems[index].Id)
				s.Set(key, string(marshaledItem))
				s.SetTTL(key, cacheTTL)
			}

			s.FastForward(cacheTTL + time.Second)

			items := itemsCache.GetItemsByIds(context.TODO(), []int{1, 2, 3})
			if len(items) != 0 {
				object, _ := json.MarshalIndent(items, "", "  ")
				fmt.Println(string(object))

				t.Error("all items should be expired")
			}
		})
	})

	t.Run("set_items_by_ids", func(t *testing.T) {
		var itemKeyPrefix = "item:%d"

		t.Run("set_no_items", func(t *testing.T) {
			s.FlushAll()

			_, err := s.Get(fmt.Sprintf(itemKeyPrefix, 1))
			if !errors.Is(err, miniredis.ErrKeyNotFound) {
				t.Errorf("expecting nil error but got something else: %v", err)
			}
		})

		t.Run("set_2_items", func(t *testing.T) {
			s.FlushAll()

			itemsCache.SetItemsByIds(context.TODO(), testItems[1:3])

			item2, err := s.Get(fmt.Sprintf(itemKeyPrefix, 2))
			if err != nil {
				t.Errorf("expecting no error but got one %v", err)
			} else if len(item2) == 0 {
				t.Error("expecting item-id 2 not be empty")
			}

			item3, err := s.Get(fmt.Sprintf(itemKeyPrefix, 3))
			if err != nil {
				t.Errorf("expecting no error but got one %v", err)
			} else if len(item3) == 0 {
				t.Error("expecting item-id 3 not be empty")
			}
		})
	})
}
