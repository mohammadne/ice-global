package entities

type Cart struct {
	Id     int
	Cookie string
	Status CartStatus
}

type CartStatus string

const (
	CartStatusOpen   CartStatus = "open"
	CartStatusClosed CartStatus = "closed"
)
