package storage_test

import (
	"testing"

	"github.com/mohammadne/ice-global/internal/config"
	"github.com/mohammadne/ice-global/pkg/mysql"
	"github.com/mohammadne/ice-global/pkg/redis"
)

var (
	mysqlDatabase *mysql.Mysql
	redisDatabase *redis.Redis
)

func TestMain(t *testing.M) {
	config, err := config.LoadDefaults(true, "/../..")
	if err != nil {
		panic(err)
	}

	mysqlDatabase, err = mysql.Open(config.Mysql)
	if err != nil {
		panic(err)
	}

	t.Run()
}
