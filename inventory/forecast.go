package inventory

import (
	"errors"
	"sort"
)

type DemandForecast struct {
	ProductID    string
	Current      int64
	AverageIssue int64
	WeeksCovered int64
	ReorderPoint int64
}

func (s *Service) Forecast(productID string, weeks int64) (DemandForecast, error) {
	if weeks <= 0 {
		return DemandForecast{}, errors.New("weeks must be positive")
	}
	product, err := s.products.Get(productID)
	if err != nil {
		return DemandForecast{}, err
	}
	movements, err := s.ledger.ForProduct(productID)
	if err != nil {
		return DemandForecast{}, err
	}
	var issued, count int64
	for _, movement := range movements {
		if movement.Kind == Issue {
			issued += movement.Quantity
			count++
		}
	}
	average := int64(0)
	if count > 0 {
		average = (issued + count - 1) / count
	}
	return DemandForecast{ProductID: productID, Current: product.Stock, AverageIssue: average, WeeksCovered: weeks, ReorderPoint: average * weeks}, nil
}

func (s *Service) ForecastAll(weeks int64) ([]DemandForecast, error) {
	products, err := s.products.ActiveProducts()
	if err != nil {
		return nil, err
	}
	result := make([]DemandForecast, 0, len(products))
	for _, product := range products {
		forecast, err := s.Forecast(product.ID, weeks)
		if err != nil {
			return nil, err
		}
		result = append(result, forecast)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ProductID < result[j].ProductID })
	return result, nil
}

func (s *Service) NeedsReorder(forecast DemandForecast) bool {
	return forecast.Current <= forecast.ReorderPoint
}
