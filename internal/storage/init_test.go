package storage_test

import (
	"testing"

	"github.com/mohammadne/ice-global/internal/config"
	"github.com/mohammadne/ice-global/pkg/mysql"
)

var (
	database *mysql.Mysql
)

func TestMain(t *testing.M) {
	config, err := config.LoadDefaults(true, "/../..")
	if err != nil {
		panic(err)
	}

	database, err = mysql.Open(config.Mysql)
	if err != nil {
		panic(err)
	}

	t.Run()
}
