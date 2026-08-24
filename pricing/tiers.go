package pricing

import (
	"errors"
	"sort"
)

type PriceTier struct {
	MinimumQuantity int64
	DiscountPercent int64
}

func ValidateTiers(tiers []PriceTier) error {
	if len(tiers) == 0 {
		return errors.New("at least one price tier is required")
	}
	ordered := append([]PriceTier(nil), tiers...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].MinimumQuantity < ordered[j].MinimumQuantity })
	previous := int64(0)
	for _, tier := range ordered {
		if tier.MinimumQuantity <= previous || tier.DiscountPercent < 0 || tier.DiscountPercent > 100 {
			return errors.New("price tier is invalid")
		}
		previous = tier.MinimumQuantity
	}
	return nil
}

func TierForQuantity(tiers []PriceTier, quantity int64) (PriceTier, error) {
	if err := ValidateTiers(tiers); err != nil {
		return PriceTier{}, err
	}
	selected := tiers[0]
	for _, tier := range tiers {
		if quantity >= tier.MinimumQuantity {
			selected = tier
		}
	}
	return selected, nil
}

func BulkQuote(cents, quantity int64, tiers []PriceTier) (int64, error) {
	if quantity <= 0 {
		return 0, errors.New("quantity must be positive")
	}
	tier, err := TierForQuantity(tiers, quantity)
	if err != nil {
		return 0, err
	}
	discounted, err := Discount(cents, tier.DiscountPercent)
	if err != nil {
		return 0, err
	}
	return discounted * quantity, nil
}
