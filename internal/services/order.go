package services

import (
	"encoding/json"
	"github.com/IlyaZayats/servord/internal/entities"
	"github.com/IlyaZayats/servord/internal/interfaces"
)

type OrderService struct {
	repo interfaces.OrderRepository
}

func NewOrderService(repo interfaces.OrderRepository) (*OrderService, error) {
	return &OrderService{repo: repo}, nil
}

func (s *OrderService) InitCache() ([]string, [][]byte, error) {
	orders, err := s.repo.GetOrders()
	if err != nil {
		return nil, nil, err
	}
	var marshalled [][]byte
	var keys []string
	for _, order := range orders {
		keys = append(keys, order.OrderUid)
		tmp, err := json.Marshal(order)
		marshalled = append(marshalled, tmp)
		if err != nil {
			return nil, nil, err
		}
	}
	return keys, marshalled, nil
}

func (s *OrderService) InsertOrder(order []byte) (string, error) {
	var orderEntity entities.Order
	err := json.Unmarshal(order, &orderEntity)
	if err != nil {
		return "", err
	}
	if err := s.repo.InsertOrder(orderEntity); err != nil {
		return "", err
	}
	return orderEntity.OrderUid, nil
}
