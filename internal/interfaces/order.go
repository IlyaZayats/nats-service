package interfaces

import "github.com/IlyaZayats/servord/internal/entities"

type OrderRepository interface {
	GetOrders() ([]entities.Order, error)
	InsertOrder(order entities.Order) error
}
