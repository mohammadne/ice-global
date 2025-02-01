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
	CreateCartItem(ctx context.Context, cartItem *CartItem) (id int, err error)
	AllCartItemsByCartId(ctx context.Context, cartId int) ([]CartItem, error)
	RetrieveCartItemByCartIdAndItemId(ctx context.Context, cartId, itemId int) (*CartItem, error)
	UpdateCartItem(ctx context.Context, cartItem *CartItem) error
	DeleteCartItemById(ctx context.Context, id int) error
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

func (ci *cartItems) CreateCartItem(ctx context.Context, cartItem *CartItem) (id int, err error) {
	query := `
	INSERT INTO cart_items (cart_id, item_id, quantity, created_at)
	VALUES (?, ?, ?, ?)`

	result, err := ci.database.ExecContext(ctx, query,
		cartItem.CartId, cartItem.ItemId, cartItem.Quantity, cartItem.CreatedAt,
	)
	if err != nil {
		return -1, fmt.Errorf("error insert cart-item into database: %v", err)
	}
	id64, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("error retrieving last inserted id: %v", err)
	}
	id = int(id64)

	return id, nil
}

func (ci *cartItems) AllCartItemsByCartId(ctx context.Context, cartId int) (result []CartItem, err error) {
	query := `
	SELECT id, cart_id, item_id, quantity, created_at, updated_at, deleted_at
	FROM cart_items
	WHERE cart_id = ?`

	rows, err := ci.database.QueryContext(ctx, query, cartId)
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

func (ci *cartItems) RetrieveCartItemByCartIdAndItemId(ctx context.Context, cartId, itemId int) (result *CartItem, err error) {
	query := `
	SELECT id, cart_id, item_id, quantity, created_at, updated_at, deleted_at
	FROM cart_items
	WHERE cart_id = ? AND item_id = ?`

	result = &CartItem{}
	err = ci.database.QueryRowContext(ctx, query, cartId, itemId).Scan(
		&result.Id, &result.CartId, &result.ItemId, &result.Quantity, &result.CreatedAt, &result.UpdatedAt, &result.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorCartItemNotFound
		}
		return nil, fmt.Errorf("error retrieving cart-item: %v", err)
	}

	return result, nil
}

var (
	ErrorCartItemNotUpdated = errors.New("cart-item has not been updated")
)

func (ci *cartItems) UpdateCartItem(ctx context.Context, cartItem *CartItem) error {
	query := `
	UPDATE cart_items SET cart_id = ?, item_id = ?, quantity = ?, updated_at = ?
	WHERE id = ?`

	result, err := ci.database.ExecContext(ctx, query,
		cartItem.CartId, cartItem.ItemId, cartItem.Quantity, cartItem.UpdatedAt, cartItem.Id,
	)
	if err != nil {
		return fmt.Errorf("error update cart-item into database: %v", err)
	}

	if counts, err := result.RowsAffected(); err != nil {
		return err
	} else if counts == 0 {
		return ErrorCartItemNotUpdated
	}

	return nil
}

func (ci *cartItems) DeleteCartItemById(ctx context.Context, id int) error {
	query := `
	UPDATE cart_items SET deleted_at = ?
	WHERE id = ?`

	result, err := ci.database.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error deleting cart-item from database: %v", err)
	}

	if counts, err := result.RowsAffected(); err != nil {
		return err
	} else if counts == 0 {
		return errors.New("cart-item has not been deleted")
	}

	return nil
}
