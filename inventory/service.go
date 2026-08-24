package inventory

import (
	"errors"
	"fmt"
	"strconv"

	"campgoods/catalog"
)

type Service struct {
	products *catalog.Service
	ledger   *Ledger
}

func NewService(products *catalog.Service, ledger *Ledger) *Service {
	return &Service{products: products, ledger: ledger}
}

func (s *Service) Record(id, productID string, kind MovementType, quantity int64, reason string) (InventoryMovement, error) {
	if quantity == 0 {
		return InventoryMovement{}, errors.New("quantity must not be zero")
	}
	product, err := s.products.Get(productID)
	if err != nil {
		return InventoryMovement{}, err
	}
	change := quantity
	if kind == Issue {
		change = -quantity
	}
	if kind == Adjust && quantity < 0 && product.Stock+quantity < 0 {
		return InventoryMovement{}, errors.New("adjustment would make stock negative")
	}
	newBalance := product.Stock + change
	if newBalance < 0 {
		return InventoryMovement{}, fmt.Errorf("insufficient stock: %d", product.Stock)
	}
	movement, err := NewMovement(id, productID, kind, quantity, reason, newBalance)
	if err != nil {
		return InventoryMovement{}, err
	}
	product.Stock = newBalance
	product.Version++
	if err := s.productsRepoSave(product); err != nil {
		return InventoryMovement{}, err
	}
	if err := s.ledger.Append(movement); err != nil {
		return InventoryMovement{}, err
	}
	return movement, nil
}

func (s *Service) productsRepoSave(product catalog.Product) error {
	return s.products.SaveForInventory(product)
}

func (s *Service) CurrentStock(productID string) (int64, error) {
	product, err := s.products.Get(productID)
	if err != nil {
		return 0, err
	}
	return product.Stock, nil
}

func (s *Service) History(productID string) ([]InventoryMovement, error) {
	return s.ledger.ForProduct(productID)
}

func (s *Service) Summary(productID string) (string, error) {
	stock, err := s.CurrentStock(productID)
	if err != nil {
		return "", err
	}
	return productID + "=" + strconv.FormatInt(stock, 10), nil
}
