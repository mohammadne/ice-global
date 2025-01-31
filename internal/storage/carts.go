package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mohammadne/ice-global/pkg/mysql"
)

type Cart struct {
	Id        int
	UserId    int
	Status    string
	CreatedAt time.Time
	DeletedAt sql.NullTime
}

type Carts interface {
	CreateCart(ctx context.Context, user *User) (id int, err error)
	RetrieveCartByUserIdAndStatus(ctx context.Context, userId int, status string) (result *Cart, err error)
}

func NewCarts(mysql *mysql.Mysql) Carts {
	return &carts{
		database: mysql,
	}
}

type carts struct {
	database *mysql.Mysql
}

func (c *carts) CreateCart(ctx context.Context, user *User) (id int, err error) {
	query := `
	INSERT INTO carts (user_id, status, created_at)
	VALUES (:user_id, :status, :created_at)
	RETURNING id INTO :id_out`

	_, err = c.database.ExecContext(ctx, query,
		sql.Named("user_id", user.Cookie),
		sql.Named("status", user.Cookie),
		sql.Named("created_at", user.CreatedAt),
		sql.Named("id_out", sql.Out{Dest: &id}),
	)

	if err != nil {
		return -1, fmt.Errorf("error insert cart into database: %v", err)
	}

	return id, nil
}

var (
	ErrorCartNotFound = errors.New("cart not found")
)

func (c *carts) RetrieveCartByUserIdAndStatus(ctx context.Context, userId int, status string) (result *Cart, err error) {
	query := `
	SELECT id, user_id, status, created_at, deleted_at
	FROM carts
	WHERE user_id = :user_id AND status = :status`

	result = &Cart{}
	err = c.database.QueryRowContext(ctx, query, sql.Named("user_id", userId), sql.Named("status", status)).Scan(
		&result.Id, &result.UserId, &result.Status, &result.CreatedAt, &result.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorCartNotFound
		}
		return nil, fmt.Errorf("error retrieing user: %v", err)
	}

	return result, nil
}
