package reporting

import (
	"sort"

	"campgoods/catalog"
)

type CatalogMetrics struct {
	TotalProducts    int
	ActiveProducts   int
	SellableProducts int
	UnitsInStock     int64
	Categories       map[string]int
}

func (s *Service) Metrics() (CatalogMetrics, error) {
	products, err := s.products.ActiveProducts()
	if err != nil {
		return CatalogMetrics{}, err
	}
	all, err := s.products.List(catalog.ProductFilter{}, 1, 100000)
	if err != nil {
		return CatalogMetrics{}, err
	}
	metrics := CatalogMetrics{TotalProducts: all.Total, ActiveProducts: len(products), Categories: make(map[string]int)}
	for _, product := range all.Items {
		metrics.UnitsInStock += product.Stock
		metrics.Categories[product.Category]++
		if product.IsSellable() {
			metrics.SellableProducts++
		}
	}
	return metrics, nil
}

func (m CatalogMetrics) CategoryNames() []string {
	names := make([]string, 0, len(m.Categories))
	for category := range m.Categories {
		names = append(names, category)
	}
	sort.Strings(names)
	return names
}

func (m CatalogMetrics) StockStatus() string {
	if m.TotalProducts == 0 {
		return "empty"
	}
	if m.SellableProducts == 0 {
		return "unavailable"
	}
	return "operational"
}

func (m CatalogMetrics) CategoryCount(category string) int {
	return m.Categories[category]
}
