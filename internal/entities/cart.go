package entities

import "time"

type Cart struct {
	Id        int
	UserId    int
	Status    CartStatus
	CreatedAt time.Time
	DeletedAt time.Time
}

type CartStatus string

const (
	CartStatusOpen   = "open"
	CartStatusClosed = "closed"
)
