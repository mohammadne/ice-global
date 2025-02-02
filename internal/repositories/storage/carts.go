package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mohammadne/ice-global/internal/entities"
	"github.com/mohammadne/ice-global/pkg/mysql"
)

type Cart struct {
	Id        int
	Cookie    string
	Status    string
	CreatedAt time.Time
	DeletedAt sql.NullTime
}

type Carts interface {
	CreateCart(ctx context.Context, cart *Cart) (id int, err error)
	RetrieveCartById(ctx context.Context, id int) (*Cart, error)
	RetrieveCartByCookieAndStatus(ctx context.Context, cookie string, status entities.CartStatus) (*Cart, error)
}

func NewCarts(mysql *mysql.Mysql) Carts {
	return &carts{
		database: mysql,
	}
}

type carts struct {
	database *mysql.Mysql
}

func (c *carts) CreateCart(ctx context.Context, cart *Cart) (id int, err error) {
	query := `
	INSERT INTO cart_entities (session_id, status, created_at)
	VALUES (?, ?, ?)`

	result, err := c.database.ExecContext(ctx, query,
		cart.Cookie, cart.Status, cart.CreatedAt,
	)
	if err != nil {
		return -1, fmt.Errorf("error insert cart into database: %v", err)
	}
	id64, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("error retrieving last inserted id: %v", err)
	}
	id = int(id64)

	return id, nil
}

var (
	ErrorCartNotFound = errors.New("cart not found")
)

func (c *carts) RetrieveCartById(ctx context.Context, id int) (result *Cart, err error) {
	query := `
	SELECT id, session_id, status, created_at, deleted_at
	FROM cart_entities
	WHERE id = ?`

	result = &Cart{}
	err = c.database.QueryRowContext(ctx, query, id).Scan(
		&result.Id, &result.Cookie, &result.Status, &result.CreatedAt, &result.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorCartNotFound
		}
		return nil, fmt.Errorf("error retrieving cart: %v", err)
	}

	return result, nil
}

func (c *carts) RetrieveCartByCookieAndStatus(ctx context.Context, cookie string, status entities.CartStatus,
) (result *Cart, err error) {
	query := `
	SELECT id, session_id, status, created_at, deleted_at
	FROM cart_entities
	WHERE session_id = ? AND status = ?`

	result = &Cart{}
	err = c.database.QueryRowContext(ctx, query, cookie, status).Scan(
		&result.Id, &result.Cookie, &result.Status, &result.CreatedAt, &result.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorCartNotFound
		}
		return nil, fmt.Errorf("error retrieving cart: %v", err)
	}

	return result, nil
}
