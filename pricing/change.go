package pricing

import (
	"errors"
	"strings"
)

type PriceChange struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Before    int64  `json:"before_cents"`
	After     int64  `json:"after_cents"`
	Reason    string `json:"reason"`
	Approved  bool   `json:"approved"`
}

func (c PriceChange) Validate() error {
	if strings.TrimSpace(c.ID) == "" || strings.TrimSpace(c.ProductID) == "" {
		return errors.New("price change identity is required")
	}
	if c.Before <= 0 || c.After <= 0 {
		return errors.New("price values must be positive")
	}
	if c.Before == c.After {
		return errors.New("price must change")
	}
	if strings.TrimSpace(c.Reason) == "" {
		return errors.New("price change reason is required")
	}
	return nil
}

func NewChange(id, productID string, before, after int64, reason string) (PriceChange, error) {
	c := PriceChange{ID: strings.TrimSpace(id), ProductID: strings.TrimSpace(productID), Before: before, After: after, Reason: strings.TrimSpace(reason)}
	if err := c.Validate(); err != nil {
		return PriceChange{}, err
	}
	return c, nil
}

func (c PriceChange) Delta() int64 {
	return c.After - c.Before
}

func (c PriceChange) Direction() string {
	if c.Delta() > 0 {
		return "increase"
	}
	return "decrease"
}
