package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/mohammadne/ice-global/pkg/metrics"
)

type Mysql struct {
	*sqlx.DB
	migrations string
	Vectors    *vectors
}

type vectors struct {
	Counter   metrics.Counter
	Histogram metrics.Histogram
}

const (
	driver      = "mysql"
	pingTimeout = time.Second * 20
)

func Open(cfg *Config, migrations, namespace, subsystem string) (*Mysql, error) {
	connString := fmt.Sprintf("%s:%s@(%s:%d)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

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

	var vectors vectors
	name := "mysql"

	counterLabels := []string{"table, function", "status"}
	vectors.Counter, err = metrics.RegisterCounter(name, namespace, subsystem, counterLabels)
	if err != nil {
		return nil, fmt.Errorf("error while registering counter vector: %v", err)
	}

	histogramLabels := []string{"table, function"}
	vectors.Histogram, err = metrics.RegisterHistogram(name, namespace, subsystem, histogramLabels)
	if err != nil {
		return nil, fmt.Errorf("error while registering histogram vector: %v", err)
	}

	r := &Mysql{DB: database, migrations: migrations, Vectors: &vectors}

	return r, nil
}

func (m *Mysql) MigrateUp(ctx context.Context) error {
	migrator := func(m *migrate.Migrate) error { return m.Up() }
	return m.migrate(m.migrations, migrator)
}

func (m *Mysql) MigrateDown(ctx context.Context) error {
	migrator := func(m *migrate.Migrate) error { return m.Down() }
	return m.migrate(m.migrations, migrator)
}

func (m *Mysql) migrate(source string, migrator func(*migrate.Migrate) error) error {
	instance, err := mysql.WithInstance(m.DB.DB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("error creating migrate instance\n%v", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(source, driver, instance)
	if err != nil {
		return fmt.Errorf("error loading migration files\n%v", err)
	}

	if err := migrator(migration); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error doing migrations\n%v", err)
	}

	return nil
}
