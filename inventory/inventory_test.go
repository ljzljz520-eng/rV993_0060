package inventory

import (
	"path/filepath"
	"testing"

	"campgoods/catalog"
	"campgoods/store"
)

func TestInventoryRecordAndReconcile(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "inventory.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	products := catalog.NewService(catalog.NewRepository(db))
	if _, err := products.Register("p1", catalog.ProductDraft{SKU: "P1", Name: "Chair", Category: "chair", Unit: "piece", PriceCents: 1000}); err != nil {
		t.Fatal(err)
	}
	service := NewService(products, NewLedger(db))
	if _, err := service.Record("m1", "p1", Receive, 10, "delivery"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Record("m2", "p1", Issue, 3, "sale"); err != nil {
		t.Fatal(err)
	}
	reconciliation, err := service.Reconcile("p1", 8)
	if err != nil || reconciliation.Difference != 1 {
		t.Fatalf("reconciliation=%+v err=%v", reconciliation, err)
	}
	if _, err := service.ApplyReconciliation("m3", reconciliation, "count"); err != nil {
		t.Fatal(err)
	}
	stock, err := service.CurrentStock("p1")
	if err != nil || stock != 8 {
		t.Fatalf("stock=%d err=%v", stock, err)
	}
}
