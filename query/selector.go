package query

import (
	"errors"
	"strings"

	"campgoods/catalog"
)

type Selector struct {
	Filter     catalog.ProductFilter
	SortField  string
	Descending bool
	Page       int
	PageSize   int
}

func NewSelector() Selector {
	return Selector{Filter: catalog.ProductFilter{ActiveOnly: true}, SortField: "id", Page: 1, PageSize: 10}
}

func (s Selector) WithQuery(query string) Selector {
	s.Filter.Query = strings.TrimSpace(query)
	return s
}
func (s Selector) WithCategory(category string) Selector {
	s.Filter.Category = strings.TrimSpace(category)
	return s
}
func (s Selector) InStock() Selector              { s.Filter.InStockOnly = true; return s }
func (s Selector) OnPage(page, size int) Selector { s.Page = page; s.PageSize = size; return s }

func (s Selector) Validate() error {
	if s.Page < 1 || s.PageSize < 1 || s.PageSize > 100 {
		return errors.New("selector page is invalid")
	}
	if s.SortField != "id" && s.SortField != "name" && s.SortField != "category" && s.SortField != "stock" {
		return errors.New("selector sort field is invalid")
	}
	return nil
}

func (s Selector) Apply(products []catalog.Product) ([]catalog.Product, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}
	filtered := make([]catalog.Product, 0, len(products))
	for _, product := range products {
		if s.Filter.Matches(product) {
			filtered = append(filtered, product)
		}
	}
	return catalog.SortProducts(filtered, s.SortField, s.Descending), nil
}

func (s Selector) Window(products []catalog.Product) ([]catalog.Product, error) {
	ordered, err := s.Apply(products)
	if err != nil {
		return nil, err
	}
	start := (s.Page - 1) * s.PageSize
	if start >= len(ordered) {
		return []catalog.Product{}, nil
	}
	end := start + s.PageSize
	if end > len(ordered) {
		end = len(ordered)
	}
	return ordered[start:end], nil
}

func (s Selector) Describe() string {
	parts := []string{"active"}
	if s.Filter.InStockOnly {
		parts = append(parts, "in-stock")
	}
	if s.Filter.Category != "" {
		parts = append(parts, "category="+s.Filter.Category)
	}
	if s.Filter.Query != "" {
		parts = append(parts, "query="+s.Filter.Query)
	}
	return strings.Join(parts, ",")
}
