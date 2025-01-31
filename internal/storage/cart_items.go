package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mohammadne/ice-global/pkg/mysql"
)

type CartItem struct {
	Id        int
	CartId    int
	ItemId    int
	Quantity  int
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
}

type CartItems interface {
	AllCartItemsByCartId(ctx context.Context, cartId int) (result []CartItem, err error)
}

func NewCartItems(mysql *mysql.Mysql) CartItems {
	return &cartItems{
		database: mysql,
	}
}

type cartItems struct {
	database *mysql.Mysql
}

var (
	ErrorCartItemNotFound = errors.New("cart-item(s) not found")
)

func (ci *cartItems) AllCartItemsByCartId(ctx context.Context, cartId int) (result []CartItem, err error) {
	query := `
	SELECT id, cart_id, item_id, quantity, created_at, updated_at, deleted_at
	FROM cart_items
	WHERE cart_id = :cart_id`

	rows, err := ci.database.QueryContext(ctx, query, sql.Named("cart_id", cartId))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorCartItemNotFound
		}
		return nil, fmt.Errorf("error query cart-items: %v", err)
	}
	defer rows.Close() // ignore error

	result = make([]CartItem, 0)
	for rows.Next() {
		item := CartItem{}
		err = rows.Scan(&item.Id, &item.CartId, &item.ItemId, &item.Quantity,
			&item.CreatedAt, &item.UpdatedAt, &item.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning cart-item result row: %v", err)
		}
		result = append(result, item)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error scanning cart-items result rows: %v", err)
	}

	return result, nil
}
