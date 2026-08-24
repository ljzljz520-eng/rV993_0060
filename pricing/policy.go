package pricing

import (
	"errors"
	"fmt"
	"strings"

	"campgoods/catalog"
)

type PricePolicy struct {
	MinimumCents       int64
	MaximumCents       int64
	MaxIncreasePercent int64
	MaxDecreasePercent int64
}

func (p PricePolicy) Validate() error {
	if p.MinimumCents <= 0 || p.MaximumCents < p.MinimumCents {
		return errors.New("price bounds are invalid")
	}
	if p.MaxIncreasePercent < 0 || p.MaxDecreasePercent < 0 {
		return errors.New("price percentages cannot be negative")
	}
	return nil
}

func (p PricePolicy) Check(before, after int64) error {
	if err := p.Validate(); err != nil {
		return err
	}
	if after < p.MinimumCents || after > p.MaximumCents {
		return fmt.Errorf("price %d is outside policy bounds", after)
	}
	if before <= 0 {
		return errors.New("current price is invalid")
	}
	delta := after - before
	if delta > 0 && delta*100 > before*p.MaxIncreasePercent {
		return errors.New("price increase exceeds policy")
	}
	if delta < 0 && (-delta)*100 > before*p.MaxDecreasePercent {
		return errors.New("price decrease exceeds policy")
	}
	return nil
}

func (s *Service) ChangeWithPolicy(id, productID string, after int64, reason string, policy PricePolicy) (PriceChange, error) {
	product, err := s.products.Get(productID)
	if err != nil {
		return PriceChange{}, err
	}
	if err := policy.Check(product.PriceCents, after); err != nil {
		return PriceChange{}, err
	}
	return s.Change(id, productID, after, reason, true)
}

func (s *Service) Explain(change PriceChange) string {
	return strings.TrimSpace(fmt.Sprintf("%s: %s by %d cents", change.ProductID, change.Direction(), change.Delta()))
}

func CompareProducts(left, right catalog.Product) int {
	if left.PriceCents < right.PriceCents {
		return -1
	}
	if left.PriceCents > right.PriceCents {
		return 1
	}
	if left.ID < right.ID {
		return -1
	}
	if left.ID > right.ID {
		return 1
	}
	return 0
}
