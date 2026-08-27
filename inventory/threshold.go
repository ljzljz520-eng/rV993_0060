package inventory

import (
	"sort"

	"campgoods/catalog"
)

type StockAlert struct {
	ProductID string
	Name      string
	Stock     int64
	Minimum   int64
	Severity  string
}

func (s *Service) Alerts(minimum int64) ([]StockAlert, error) {
	products, err := s.products.ActiveProducts()
	if err != nil {
		return nil, err
	}
	alerts := make([]StockAlert, 0)
	for _, product := range products {
		if product.Stock <= minimum {
			severity := "watch"
			if product.Stock == 0 {
				severity = "critical"
			}
			alerts = append(alerts, StockAlert{ProductID: product.ID, Name: product.Name, Stock: product.Stock, Minimum: minimum, Severity: severity})
		}
	}
	sort.Slice(alerts, func(i, j int) bool {
		if alerts[i].Stock != alerts[j].Stock {
			return alerts[i].Stock < alerts[j].Stock
		}
		return alerts[i].ProductID < alerts[j].ProductID
	})
	return alerts, nil
}

func (s *Service) Availability(product catalog.Product) string {
	if !product.Active {
		return "archived"
	}
	if product.Stock == 0 {
		return "out-of-stock"
	}
	if product.Stock < 3 {
		return "low-stock"
	}
	return "available"
}

func (s *Service) CanIssue(productID string, quantity int64) (bool, error) {
	if quantity <= 0 {
		return false, nil
	}
	product, err := s.products.Get(productID)
	if err != nil {
		return false, err
	}
	return product.Active && product.Stock >= quantity, nil
}
