package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mohammadne/ice-global/pkg/mysql"
)

type Item struct {
	Id        int
	Name      string
	Price     int
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type Items interface {
	AllItems(ctx context.Context) ([]Item, error)
}

func NewItems(mysql *mysql.Mysql) Items {
	return &items{
		database: mysql,
	}
}

type items struct {
	database *mysql.Mysql
}

func (i *items) AllItems(ctx context.Context) (result []Item, err error) {
	query := `
	SELECT id, name, price, created_at, updated_at
	FROM items`

	rows, err := i.database.QueryContext(ctx, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no rows for items")
		}
		return nil, fmt.Errorf("error query items: %v", err)
	}
	defer rows.Close() // ignore error

	result = make([]Item, 0)
	for rows.Next() {
		item := Item{}
		err = rows.Scan(&item.Id, &item.Name, &item.Price, &item.CreatedAt, &item.UpdatedAt)
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
