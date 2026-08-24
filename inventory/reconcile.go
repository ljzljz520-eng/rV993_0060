package inventory

import (
	"errors"
	"fmt"
	"sort"
)

type Reconciliation struct {
	ProductID   string
	Recorded    int64
	Expected    int64
	Difference  int64
	NeedsAdjust bool
}

func (s *Service) Reconcile(productID string, expected int64) (Reconciliation, error) {
	if expected < 0 {
		return Reconciliation{}, errors.New("expected stock cannot be negative")
	}
	product, err := s.products.Get(productID)
	if err != nil {
		return Reconciliation{}, err
	}
	difference := expected - product.Stock
	return Reconciliation{ProductID: productID, Recorded: product.Stock, Expected: expected, Difference: difference, NeedsAdjust: difference != 0}, nil
}

func (s *Service) ApplyReconciliation(movementID string, reconciliation Reconciliation, reason string) (InventoryMovement, error) {
	if !reconciliation.NeedsAdjust {
		return InventoryMovement{}, errors.New("reconciliation has no difference")
	}
	if reconciliation.Expected < 0 {
		return InventoryMovement{}, errors.New("expected stock cannot be negative")
	}
	product, err := s.products.Get(reconciliation.ProductID)
	if err != nil {
		return InventoryMovement{}, err
	}
	if product.Stock+reconciliation.Difference != reconciliation.Expected {
		return InventoryMovement{}, errors.New("reconciliation is stale")
	}
	movement, err := NewMovement(movementID, reconciliation.ProductID, Adjust, reconciliation.Difference, reason, reconciliation.Expected)
	if err != nil {
		return InventoryMovement{}, err
	}
	product.Stock = reconciliation.Expected
	product.Version++
	if err := s.products.SaveForInventory(product); err != nil {
		return InventoryMovement{}, err
	}
	if err := s.ledger.Append(movement); err != nil {
		return InventoryMovement{}, err
	}
	return movement, nil
}

type StockVariance struct {
	ProductID  string
	Difference int64
}

func (s *Service) ReconcileMany(expected map[string]int64) ([]StockVariance, error) {
	if expected == nil {
		return nil, errors.New("expected stock map is required")
	}
	keys := make([]string, 0, len(expected))
	for productID := range expected {
		keys = append(keys, productID)
	}
	sort.Strings(keys)
	variances := make([]StockVariance, 0)
	for _, productID := range keys {
		result, err := s.Reconcile(productID, expected[productID])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", productID, err)
		}
		if result.NeedsAdjust {
			variances = append(variances, StockVariance{ProductID: productID, Difference: result.Difference})
		}
	}
	return variances, nil
}

func (s *Service) MovementCount(productID string) (int, error) {
	movements, err := s.ledger.ForProduct(productID)
	return len(movements), err
}
