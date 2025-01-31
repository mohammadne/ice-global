package entities

import "time"

type CartItem struct {
	Id        int
	CartId    int
	ItemId    int
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
