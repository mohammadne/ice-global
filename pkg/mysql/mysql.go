package mysql

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

type Mysql struct {
	*sqlx.DB
	migrations string
}

const (
	driver      = "mysql"
	pingTimeout = time.Second * 20
)

func Open(cfg *Config, migrations string) (*Mysql, error) {
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

	r := &Mysql{DB: database, migrations: migrations}

	return r, nil
}

type MigrateDirection string

const (
	MigrateDirectionUp   MigrateDirection = "UP"
	MigrateDirectionDown MigrateDirection = "DOWN"
)

func (m *Mysql) Migrate(direction MigrateDirection) error {
	if m.migrations == "" {
		return fmt.Errorf("no migration directory has been given")
	}

	instance, err := mysql.WithInstance(m.DB.DB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("error creating migrate instance\n%v", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(m.migrations, driver, instance)
	if err != nil {
		return fmt.Errorf("error loading migration files\n%v", err)
	}

	var migrator func(m *migrate.Migrate) error

	switch direction {
	case MigrateDirectionUp:
		migrator = func(m *migrate.Migrate) error { return m.Up() }
	case MigrateDirectionDown:
		migrator = func(m *migrate.Migrate) error { return m.Down() }
	default:
		return fmt.Errorf("invalid direction has been given\n%s", direction)
	}

	if err := migrator(migration); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error doing migrations\n%v", err)
	}

	return nil
}
