package catalog

import (
	"errors"
	"fmt"
	"strings"
)

type Product struct {
	ID          string `json:"id"`
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Unit        string `json:"unit"`
	PriceCents  int64  `json:"price_cents"`
	Stock       int64  `json:"stock"`
	Active      bool   `json:"active"`
	Version     int64  `json:"version"`
}

type ProductDraft struct {
	SKU         string
	Name        string
	Category    string
	Description string
	Unit        string
	PriceCents  int64
}

func (d ProductDraft) Normalize() ProductDraft {
	d.SKU = strings.ToUpper(strings.TrimSpace(d.SKU))
	d.Name = strings.TrimSpace(d.Name)
	d.Category = strings.ToLower(strings.TrimSpace(d.Category))
	d.Description = strings.TrimSpace(d.Description)
	d.Unit = strings.TrimSpace(d.Unit)
	return d
}

func (d ProductDraft) Validate() error {
	n := d.Normalize()
	if n.SKU == "" || n.Name == "" {
		return errors.New("sku and name are required")
	}
	if !AllowedCategory(n.Category) {
		return fmt.Errorf("unsupported category %q", n.Category)
	}
	if n.Unit == "" {
		return errors.New("unit is required")
	}
	if n.PriceCents <= 0 {
		return errors.New("price must be positive")
	}
	return nil
}

func AllowedCategory(category string) bool {
	switch strings.ToLower(strings.TrimSpace(category)) {
	case "tarp", "chair", "cart", "mat":
		return true
	default:
		return false
	}
}

func NewProduct(id string, d ProductDraft) (Product, error) {
	d = d.Normalize()
	if err := d.Validate(); err != nil {
		return Product{}, err
	}
	if strings.TrimSpace(id) == "" {
		return Product{}, errors.New("product id is required")
	}
	return Product{ID: id, SKU: d.SKU, Name: d.Name, Category: d.Category, Description: d.Description, Unit: d.Unit, PriceCents: d.PriceCents, Active: true, Version: 1}, nil
}

func (p Product) Validate() error {
	if p.ID == "" || p.SKU == "" || p.Name == "" {
		return errors.New("product identity is incomplete")
	}
	if !AllowedCategory(p.Category) {
		return errors.New("product category is invalid")
	}
	if p.PriceCents <= 0 || p.Version < 1 {
		return errors.New("product commercial values are invalid")
	}
	if p.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}
	return nil
}

func (p Product) IsSellable() bool {
	return p.Active && p.Stock > 0 && p.PriceCents > 0
}

func (p Product) DisplayPrice() string {
	return fmt.Sprintf("%d.%02d", p.PriceCents/100, p.PriceCents%100)
}

func (p Product) Clone() Product {
	return Product{ID: p.ID, SKU: p.SKU, Name: p.Name, Category: p.Category, Description: p.Description, Unit: p.Unit, PriceCents: p.PriceCents, Stock: p.Stock, Active: p.Active, Version: p.Version}
}
