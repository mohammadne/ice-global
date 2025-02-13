package storage_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/mohammadne/shopping-cart-manager/internal/repositories/storage"
	"github.com/mohammadne/shopping-cart-manager/pkg/mysql"
)

var (
	mockDB sqlmock.Sqlmock

	// storages
	cartItemStorage storage.CartItems
	cartsStorage    storage.Carts
	itemsStorage    storage.Items
)

func TestMain(m *testing.M) {
	var err error
	var sqlDB *sql.DB

	sqlDB, mockDB, err = sqlmock.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not start sqlmock: %v\n", err)
		os.Exit(1) // Exit with a non-zero status code
	}
	defer sqlDB.Close()

	sqlxDB := sqlx.NewDb(sqlDB, "sqlmock")
	mysql := &mysql.Mysql{DB: sqlxDB}

	cartItemStorage = storage.NewCartItems(mysql)
	cartsStorage = storage.NewCarts(mysql)
	itemsStorage = storage.NewItems(mysql)

	m.Run()
}
