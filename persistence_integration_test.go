package campgoods_test

import (
	"path/filepath"
	"testing"

	"campgoods/catalog"
	"campgoods/reporting"
	"campgoods/workflow"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "persist.db")
	app, err := workflow.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := app.RegisterProduct("persist-01", catalog.ProductDraft{SKU: "PERSIST-01", Name: "Persistent tarp", Category: "tarp", Unit: "piece", PriceCents: 2500, Description: "stored"}); err != nil {
		t.Fatal(err)
	}
	if _, err := app.ReceiveStock("persist-move", "persist-01", 9, "opening"); err != nil {
		t.Fatal(err)
	}
	if _, err := app.ChangePrice("persist-price", "persist-01", 2750, "review"); err != nil {
		t.Fatal(err)
	}
	if err := app.AddProductNote(reporting.ProductNote{ID: "persist-note", ProductID: "persist-01", Text: "reopen me"}); err != nil {
		t.Fatal(err)
	}
	if err := app.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := workflow.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	product, err := reopened.Catalog.Get("persist-01")
	if err != nil || product.Stock != 9 || product.PriceCents != 2750 {
		t.Fatalf("product=%+v err=%v", product, err)
	}
	changes, err := reopened.Pricing.HistoryFor("persist-01")
	if err != nil || len(changes) != 1 {
		t.Fatalf("changes=%d err=%v", len(changes), err)
	}
	notes, err := reopened.Reporting.Notes("persist-01")
	if err != nil || len(notes) != 1 {
		t.Fatalf("notes=%d err=%v", len(notes), err)
	}
}
