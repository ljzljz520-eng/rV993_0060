package reporting

import (
	"path/filepath"
	"testing"

	"campgoods/catalog"
	"campgoods/inventory"
	"campgoods/pricing"
	"campgoods/store"
)

func TestSummaryAndExport(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "report.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	products := catalog.NewService(catalog.NewRepository(db))
	if _, err := products.Register("p1", catalog.ProductDraft{SKU: "P1", Name: "Mat", Category: "mat", Unit: "piece", PriceCents: 800}); err != nil {
		t.Fatal(err)
	}
	inv := inventory.NewService(products, inventory.NewLedger(db))
	prices := pricing.NewService(products, pricing.NewHistory(db))
	if _, err := inv.Record("m1", "p1", inventory.Receive, 5, "delivery"); err != nil {
		t.Fatal(err)
	}
	service := NewService(products, inv, prices, db)
	if err := service.AddNote(ProductNote{ID: "n1", ProductID: "p1", Text: "foam"}); err != nil {
		t.Fatal(err)
	}
	bundle, err := service.Export([]string{"p1"})
	if err != nil || len(bundle.Summaries) != 1 {
		t.Fatalf("bundle=%+v err=%v", bundle, err)
	}
	data, err := MarshalExport(bundle)
	if err != nil || len(data) == 0 {
		t.Fatalf("data=%q err=%v", data, err)
	}
}
