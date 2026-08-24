package pricing

import "testing"

func TestPolicyAndTierQuotes(t *testing.T) {
	policy := PricePolicy{MinimumCents: 100, MaximumCents: 10000, MaxIncreasePercent: 20, MaxDecreasePercent: 30}
	if err := policy.Check(1000, 1150); err != nil {
		t.Fatal(err)
	}
	if err := policy.Check(1000, 1301); err == nil {
		t.Fatal("expected increase rejection")
	}
	quote, err := BulkQuote(1000, 12, []PriceTier{{MinimumQuantity: 1, DiscountPercent: 0}, {MinimumQuantity: 10, DiscountPercent: 10}})
	if err != nil || quote != 10800 {
		t.Fatalf("quote=%d err=%v", quote, err)
	}
}
