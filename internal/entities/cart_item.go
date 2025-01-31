package entities

type CartItem struct {
	Id        int
	Cart      *Cart
	Item      *Item
	Quantity  int
	IsDeleted bool
}
