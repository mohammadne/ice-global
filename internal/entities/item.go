package entities

import "time"

type Item struct {
	Id        int
	Name      string
	Price     int
	CreatedAt time.Time
}
