package mysql

import (
	"context"
	"fmt"
	"io/fs"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Mysql struct {
	*sqlx.DB
}

const (
	driver      = "mysql"
	pingTimeout = time.Second * 20
)

func Open(cfg *Config) (*Mysql, error) {
	connString := fmt.Sprintf(
		"%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)

	database, err := sqlx.Open(driver, connString)
	if err != nil {
		return nil, fmt.Errorf("error while opening connection to mysql: %v", err)
	}
	database.SetMaxIdleConns(0)

	ctx, cf := context.WithTimeout(context.Background(), pingTimeout)
	defer cf()
	if err = database.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error while pinging database: %v", err)
	}

	r := &Mysql{DB: database}

	return r, nil
}

type MigrateDirection string

const (
	MigrateDirectionUp   MigrateDirection = "UP"
	MigrateDirectionDown MigrateDirection = "DOWN"
)

func (m *Mysql) Migrate(migrations []fs.DirEntry, direction MigrateDirection) error {
	if len(migrations) == 0 {
		return fmt.Errorf("no migration files has been given")
	}

	return nil
}
