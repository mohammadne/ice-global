package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mohammadne/shopping-cart-manager/pkg/mysql"
)

type Item struct {
	ID        int
	Name      string
	Price     int
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type Items interface {
	AllItems(ctx context.Context) ([]Item, error)
	AllItemsByItemIds(ctx context.Context, ids []int) ([]Item, error)
}

func NewItems(mysql *mysql.Mysql) Items {
	return &items{mysql: mysql}
}

type items struct {
	mysql *mysql.Mysql
}

var (
	ErrorItemNotFound = errors.New("item(s) not found")
)

func (i *items) AllItems(ctx context.Context) (result []Item, err error) {
	query := `
	SELECT id, name, price, created_at, updated_at
	FROM items`

	rows, err := i.mysql.QueryContext(ctx, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorItemNotFound
		}
		return nil, fmt.Errorf("error query items: %v", err)
	}
	defer rows.Close() // ignore error

	result = make([]Item, 0)
	for rows.Next() {
		item := Item{}
		err = rows.Scan(&item.ID, &item.Name, &item.Price, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning item result row: %v", err)
		}
		result = append(result, item)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error scanning items result rows: %v", err)
	}

	return result, nil
}

func (i *items) AllItemsByItemIds(ctx context.Context, ids []int) (result []Item, err error) {
	if len(ids) == 0 {
		return []Item{}, nil
	}

	query := `
	SELECT id, name, price, created_at, updated_at
	FROM items
	WHERE id IN (?)`

	expandedQuery, args, err := sqlx.In(query, ids)
	if err != nil {
		log.Fatalf("Error preparing query: %v", err)
	}
	expandedQuery = i.mysql.Rebind(expandedQuery)

	rows, err := i.mysql.QueryContext(ctx, expandedQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorItemNotFound
		}
		return nil, fmt.Errorf("error query items: %v", err)
	}
	defer rows.Close() // ignore error

	result = make([]Item, 0)
	for rows.Next() {
		item := Item{}
		err = rows.Scan(&item.ID, &item.Name, &item.Price, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning item result row: %v", err)
		}
		result = append(result, item)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error scanning items result rows: %v", err)
	}

	return result, nil
}
