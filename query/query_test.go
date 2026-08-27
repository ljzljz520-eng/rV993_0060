package query

import (
	"testing"

	"campgoods/catalog"
)

func TestSelectorWindowAndColumns(t *testing.T) {
	selector := NewSelector().WithCategory("chair").InStock().OnPage(1, 2)
	products := []catalog.Product{{ID: "a", Name: "A", Category: "chair", Stock: 1, Active: true}, {ID: "b", Name: "B", Category: "chair", Stock: 0, Active: true}, {ID: "c", Name: "C", Category: "tarp", Stock: 4, Active: true}}
	window, err := selector.Window(products)
	if err != nil || len(window) != 1 || window[0].ID != "a" {
		t.Fatalf("window=%+v err=%v", window, err)
	}
	columns, err := ParseColumns("sku,name,stock,sku")
	if err != nil || len(columns) != 3 {
		t.Fatalf("columns=%v err=%v", columns, err)
	}
	if Value(products[0], ColumnStock) != "1" {
		t.Fatal(Value(products[0], ColumnStock))
	}
}
