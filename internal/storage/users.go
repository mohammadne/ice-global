package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mohammadne/ice-global/pkg/mysql"
)

type User struct {
	Id        int
	Cookie    string
	CreatedAt time.Time
}

type Users interface {
	CreateUser(ctx context.Context, user *User) (int, error)
	RetrieveUserByCookie(ctx context.Context, cookie string) (*User, error)
}

func NewUsers(mysql *mysql.Mysql) Users {
	return &users{
		database: mysql,
	}
}

type users struct {
	database *mysql.Mysql
}

func (u *users) CreateUser(ctx context.Context, user *User) (id int, err error) {
	query := `
	INSERT INTO users (cookie, created_at)
	VALUES (:cookie, :created_at)
	RETURNING id INTO :id_out`

	_, err = u.database.ExecContext(ctx, query,
		sql.Named("cookie", user.Cookie),
		sql.Named("created_at", user.CreatedAt),
		sql.Named("id_out", sql.Out{Dest: &id}),
	)

	if err != nil {
		return -1, fmt.Errorf("error insert user into database: %v", err)
	}

	return id, nil
}

var (
	ErrorUserNotFound = errors.New("user not found")
)

func (u *users) RetrieveUserByCookie(ctx context.Context, cookie string) (result *User, err error) {
	query := `
	SELECT id, cookie, created_at
	FROM users
	WHERE cookie = :cookie`

	result = &User{}
	err = u.database.QueryRowContext(ctx, query, sql.Named("cookie", cookie)).Scan(
		&result.Id, &result.Cookie, &result.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorUserNotFound
		}
		return nil, fmt.Errorf("error retrieing user: %v", err)
	}

	return result, nil
}
