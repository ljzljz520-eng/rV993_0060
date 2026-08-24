package campgoods_test

import (
	"path/filepath"
	"testing"

	"campgoods/catalog"
	"campgoods/pricing"
	"campgoods/query"
	"campgoods/reporting"
	"campgoods/workflow"
)

func newTestApp(t *testing.T) *workflow.App {
	t.Helper()
	app, err := workflow.Open(filepath.Join(t.TempDir(), "campgoods.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = app.Close() })
	return app
}

func registerTestProduct(t *testing.T, app *workflow.App, id string) {
	t.Helper()
	_, err := app.RegisterProduct(id, catalog.ProductDraft{SKU: id, Name: "Shelter " + id, Category: "tarp", Unit: "piece", PriceCents: 1000, Description: "field item"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWorkflowProductOnboarding(t *testing.T) {
	app := newTestApp(t)
	registerTestProduct(t, app, "tarp-01")
	if _, err := app.ReceiveStock("move-01", "tarp-01", 6, "delivery"); err != nil {
		t.Fatal(err)
	}
	if _, err := app.ChangePrice("price-01", "tarp-01", 1200, "seasonal"); err != nil {
		t.Fatal(err)
	}
	summary, err := app.ProductSummary("tarp-01")
	if err != nil {
		t.Fatal(err)
	}
	if summary.Stock != 6 || summary.Product.PriceCents != 1200 || summary.Availability != "available" {
		t.Fatalf("unexpected summary: %+v", summary)
	}
}

func TestWorkflowStockAndPricing(t *testing.T) {
	app := newTestApp(t)
	registerTestProduct(t, app, "chair-01")
	if _, err := app.ReceiveStock("move-01", "chair-01", 10, "opening"); err != nil {
		t.Fatal(err)
	}
	if _, err := app.Inventory.Record("move-02", "chair-01", "issue", 3, "sale"); err != nil {
		t.Fatal(err)
	}
	if _, err := app.Pricing.ChangeWithPolicy("price-01", "chair-01", 1080, "minor increase", pricingPolicy()); err != nil {
		t.Fatal(err)
	}
	stock, err := app.Inventory.CurrentStock("chair-01")
	if err != nil || stock != 7 {
		t.Fatalf("stock=%d err=%v", stock, err)
	}
	quote, err := app.Pricing.Quote("chair-01", 2)
	if err != nil || quote != 2160 {
		t.Fatalf("quote=%d err=%v", quote, err)
	}
}

func pricingPolicy() pricing.PricePolicy {
	return pricing.PricePolicy{MinimumCents: 1, MaximumCents: 100000, MaxIncreasePercent: 20, MaxDecreasePercent: 50}
}

func TestWorkflowReportingAndRecovery(t *testing.T) {
	app := newTestApp(t)
	registerTestProduct(t, app, "mat-01")
	if _, err := app.ReceiveStock("move-01", "mat-01", 4, "opening"); err != nil {
		t.Fatal(err)
	}
	if err := app.AddProductNote(reporting.ProductNote{ID: "note-01", ProductID: "mat-01", Text: "keep dry", Pinned: true}); err != nil {
		t.Fatal(err)
	}
	view, err := app.ListProducts(query.ListRequest{Page: 1, PageSize: 10})
	if err != nil || len(view.Rows) != 1 {
		t.Fatalf("view rows=%d err=%v", len(view.Rows), err)
	}
	status, err := app.WorkflowStatus("mat-01")
	if err != nil || status != "ready" {
		t.Fatalf("status=%s err=%v", status, err)
	}
}
