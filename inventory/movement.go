package inventory

import (
	"errors"
	"fmt"
	"strings"
)

type MovementType string

const (
	Receive MovementType = "receive"
	Issue   MovementType = "issue"
	Adjust  MovementType = "adjust"
)

type InventoryMovement struct {
	ID        string       `json:"id"`
	ProductID string       `json:"product_id"`
	Kind      MovementType `json:"kind"`
	Quantity  int64        `json:"quantity"`
	Reason    string       `json:"reason"`
	Balance   int64        `json:"balance"`
}

func (m InventoryMovement) Validate() error {
	if strings.TrimSpace(m.ID) == "" || strings.TrimSpace(m.ProductID) == "" {
		return errors.New("movement identity is required")
	}
	if m.Quantity == 0 {
		return errors.New("movement quantity cannot be zero")
	}
	switch m.Kind {
	case Receive, Issue, Adjust:
	default:
		return fmt.Errorf("unsupported movement kind %q", m.Kind)
	}
	if strings.TrimSpace(m.Reason) == "" {
		return errors.New("movement reason is required")
	}
	return nil
}

func (m InventoryMovement) SignedQuantity() int64 {
	if m.Kind == Issue && m.Quantity > 0 {
		return -m.Quantity
	}
	if m.Kind == Receive && m.Quantity < 0 {
		return -m.Quantity
	}
	return m.Quantity
}

func NewMovement(id, productID string, kind MovementType, quantity int64, reason string, balance int64) (InventoryMovement, error) {
	m := InventoryMovement{ID: strings.TrimSpace(id), ProductID: strings.TrimSpace(productID), Kind: kind, Quantity: quantity, Reason: strings.TrimSpace(reason), Balance: balance}
	if err := m.Validate(); err != nil {
		return InventoryMovement{}, err
	}
	return m, nil
}
