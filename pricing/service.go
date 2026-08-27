package pricing

import (
	"errors"
	"fmt"

	"campgoods/catalog"
)

type Service struct {
	products *catalog.Service
	history  *History
}

func NewService(products *catalog.Service, history *History) *Service {
	return &Service{products: products, history: history}
}

func (s *Service) Change(id, productID string, after int64, reason string, approved bool) (PriceChange, error) {
	if !approved {
		return PriceChange{}, errors.New("price change requires approval")
	}
	product, err := s.products.Get(productID)
	if err != nil {
		return PriceChange{}, err
	}
	change, err := NewChange(id, productID, product.PriceCents, after, reason)
	if err != nil {
		return PriceChange{}, err
	}
	change.Approved = approved
	product.PriceCents = after
	product.Version++
	if err := s.products.SaveForInventory(product); err != nil {
		return PriceChange{}, err
	}
	if err := s.history.Append(change); err != nil {
		return PriceChange{}, err
	}
	return change, nil
}

func (s *Service) Current(productID string) (int64, error) {
	product, err := s.products.Get(productID)
	if err != nil {
		return 0, err
	}
	return product.PriceCents, nil
}

func (s *Service) HistoryFor(productID string) ([]PriceChange, error) {
	return s.history.ForProduct(productID)
}

func (s *Service) Quote(productID string, quantity int64) (int64, error) {
	if quantity <= 0 {
		return 0, errors.New("quantity must be positive")
	}
	price, err := s.Current(productID)
	if err != nil {
		return 0, err
	}
	if quantity > 1000 {
		return 0, fmt.Errorf("quantity exceeds quote limit")
	}
	return price * quantity, nil
}
