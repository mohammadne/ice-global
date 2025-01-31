package entities

type Cart struct {
	Id     int
	UserId int
	Status CartStatus
}

type CartStatus string

const (
	CartStatusOpen   = "open"
	CartStatusClosed = "closed"
)
