package redis

import (
	"time"
)

type Config struct {
	Address  string        `required:"true"`
	Username string        `required:"true"`
	Password string        `required:"true"`
	DB       int           `required:"true"`
	Timeout  time.Duration `required:"false" default:"5s"`
	PoolSize int           `required:"false" default:"10"`
}
