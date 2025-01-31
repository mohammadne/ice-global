package entities

import "time"

type User struct {
	Id        int
	Session   string
	CreatedAt time.Time
}
