package catalog

import (
	"strings"
	"testing"
)

func TestProductDraftValidationAndImport(t *testing.T) {
	input := "id,sku,name,category,unit,price,description\n" +
		"a,A-1,Ridge tarp,tarp,piece,1200,weather cover\n" +
		"b,B-1,Unknown,lantern,piece,300,invalid\n"
	result, err := ParseImport(strings.NewReader(input))
	if err != nil || result.RowsAccepted != 1 || result.RowsRejected != 1 {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	if result.Products[0].DisplayPrice() != "12.00" {
		t.Fatal(result.Products[0].DisplayPrice())
	}
}

func TestProductFilteringSorting(t *testing.T) {
	products := []Product{{ID: "2", SKU: "B", Name: "Chair", Category: "chair"}, {ID: "1", SKU: "A", Name: "Tarp", Category: "tarp", Stock: 2}}
	filtered := []Product{}
	filter := ProductFilter{Query: "tar", InStockOnly: true}
	for _, product := range products {
		if filter.Matches(product) {
			filtered = append(filtered, product)
		}
	}
	if len(filtered) != 1 || filtered[0].ID != "1" {
		t.Fatalf("filtered=%+v", filtered)
	}
	sorted := SortProducts(products, "name", false)
	if sorted[0].Name != "Chair" {
		t.Fatalf("sorted=%+v", sorted)
	}
}
