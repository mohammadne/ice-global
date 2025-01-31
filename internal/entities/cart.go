package entities

import "time"

type Cart struct {
	Id        int
	UserId    int
	Status    string
	CreatedAt time.Time
	DeletedAt time.Time
}
