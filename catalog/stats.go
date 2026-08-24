package catalog

import (
	"sort"
)

type CategoryStat struct {
	Category   string
	Products   int
	Units      int64
	ValueCents int64
}

func (s *Service) CategoryStats() ([]CategoryStat, error) {
	products, err := s.repo.All()
	if err != nil {
		return nil, err
	}
	byCategory := make(map[string]CategoryStat)
	for _, product := range products {
		stat := byCategory[product.Category]
		stat.Category = product.Category
		stat.Products++
		stat.Units += product.Stock
		stat.ValueCents += product.Stock * product.PriceCents
		byCategory[product.Category] = stat
	}
	result := make([]CategoryStat, 0, len(byCategory))
	for _, stat := range byCategory {
		result = append(result, stat)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Category < result[j].Category })
	return result, nil
}

func (s *Service) SearchBySKU(prefix string) ([]Product, error) {
	products, err := s.repo.All()
	if err != nil {
		return nil, err
	}
	result := make([]Product, 0)
	for _, product := range products {
		if len(prefix) <= len(product.SKU) && product.SKU[:len(prefix)] == prefix {
			result = append(result, product)
		}
	}
	return result, nil
}

func (s *Service) SearchByCategory(category string) ([]Product, error) {
	products, err := s.repo.All()
	if err != nil {
		return nil, err
	}
	result := make([]Product, 0)
	for _, product := range products {
		if product.Category == category {
			result = append(result, product)
		}
	}
	return result, nil
}

func (s *Service) InventoryValue() (int64, error) {
	products, err := s.repo.All()
	if err != nil {
		return 0, err
	}
	var total int64
	for _, product := range products {
		total += product.PriceCents * product.Stock
	}
	return total, nil
}

func (s *Service) CountByAvailability() (map[string]int, error) {
	products, err := s.repo.All()
	if err != nil {
		return nil, err
	}
	counts := map[string]int{"available": 0, "low-stock": 0, "out-of-stock": 0, "archived": 0}
	for _, product := range products {
		status := "out-of-stock"
		if !product.Active {
			status = "archived"
		} else if product.Stock >= 3 {
			status = "available"
		} else if product.Stock > 0 {
			status = "low-stock"
		}
		counts[status]++
	}
	return counts, nil
}

func (s *Service) HighestValue(limit int) ([]Product, error) {
	if limit < 1 {
		return nil, nil
	}
	products, err := s.repo.All()
	if err != nil {
		return nil, err
	}
	sort.SliceStable(products, func(i, j int) bool {
		left := products[i].PriceCents * products[i].Stock
		right := products[j].PriceCents * products[j].Stock
		if left == right {
			return products[i].ID < products[j].ID
		}
		return left > right
	})
	if limit > len(products) {
		limit = len(products)
	}
	return products[:limit], nil
}
