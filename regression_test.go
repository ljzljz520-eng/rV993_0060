package campgoods_test

import (
	"fmt"
	"testing"

	"campgoods/catalog"
	"campgoods/query"
)

func TestProductPageKeepsContiguousItems(t *testing.T) {
	app := newTestApp(t)
	for index := 1; index <= 30; index++ {
		id := fmt.Sprintf("product-%02d", index)
		if _, err := app.RegisterProduct(id, catalog.ProductDraft{SKU: id, Name: "Product " + id, Category: "chair", Unit: "piece", PriceCents: 1000 + int64(index), Description: "listed"}); err != nil {
			t.Fatal(err)
		}
	}
	view, err := app.ListProducts(query.ListRequest{Page: 2, PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(view.Page.Items) != 10 {
		t.Fatalf("expected ten items, got %d", len(view.Page.Items))
	}
	for index, product := range view.Page.Items {
		want := fmt.Sprintf("product-%02d", index+11)
		if product.ID != want {
			t.Fatalf("position %d: got %s want %s", index, product.ID, want)
		}
	}
}
