package workflow

import (
	"fmt"

	"campgoods/catalog"
)

type SeedSpec struct {
	ID, SKU, Name, Category, Unit, Description string
	PriceCents, Stock                          int64
}

func (a *App) Seed(specs []SeedSpec) BatchOutcome {
	drafts := make(map[string]catalog.ProductDraft, len(specs))
	for _, spec := range specs {
		drafts[spec.ID] = catalog.ProductDraft{SKU: spec.SKU, Name: spec.Name, Category: spec.Category, Unit: spec.Unit, Description: spec.Description, PriceCents: spec.PriceCents}
	}
	outcome := a.RegisterBatch(drafts)
	for _, spec := range specs {
		if _, err := a.ReceiveStock("seed-"+spec.ID, spec.ID, spec.Stock, "initial stock"); err != nil {
			outcome.Failures = append(outcome.Failures, fmt.Sprintf("%s: %v", spec.ID, err))
		} else {
			outcome.Completed++
		}
	}
	outcome.Requested += len(specs)
	return outcome
}

func DefaultSeed() []SeedSpec {
	return []SeedSpec{
		{ID: "tarp-001", SKU: "TARP-001", Name: "Ridge tarp", Category: "tarp", Unit: "piece", PriceCents: 18900, Stock: 8, Description: "silnylon shelter"},
		{ID: "chair-001", SKU: "CHAIR-001", Name: "Folding chair", Category: "chair", Unit: "piece", PriceCents: 6900, Stock: 12, Description: "aluminum frame"},
		{ID: "cart-001", SKU: "CART-001", Name: "Camp cart", Category: "cart", Unit: "piece", PriceCents: 12900, Stock: 4, Description: "all terrain cart"},
		{ID: "mat-001", SKU: "MAT-001", Name: "Dry mat", Category: "mat", Unit: "piece", PriceCents: 4900, Stock: 15, Description: "closed cell foam"},
	}
}
