package catalog

import (
	"errors"
	"sort"
)

type LifecycleEvent struct {
	ProductID string
	From      bool
	To        bool
	Version   int64
}

func (s *Service) Toggle(id string) (LifecycleEvent, error) {
	product, err := s.repo.Require(id)
	if err != nil {
		return LifecycleEvent{}, err
	}
	previous := product.Active
	product.Active = !product.Active
	product.Version++
	if err := s.repo.Save(product); err != nil {
		return LifecycleEvent{}, err
	}
	return LifecycleEvent{ProductID: id, From: previous, To: product.Active, Version: product.Version}, nil
}

func (s *Service) ActiveProducts() ([]Product, error) {
	products, err := s.repo.All()
	if err != nil {
		return nil, err
	}
	active := make([]Product, 0, len(products))
	for _, product := range products {
		if product.Active {
			active = append(active, product)
		}
	}
	sort.Slice(active, func(i, j int) bool { return active[i].Name < active[j].Name })
	return active, nil
}

func (s *Service) DeleteInactive(id string) error {
	product, err := s.repo.Require(id)
	if err != nil {
		return err
	}
	if product.Active {
		return errors.New("active product cannot be deleted")
	}
	return s.repo.Delete(id)
}

func (s *Service) ValidateCatalog() ([]string, error) {
	products, err := s.repo.All()
	if err != nil {
		return nil, err
	}
	issues := make([]string, 0)
	seenSKUs := make(map[string]string)
	for _, product := range products {
		if err := product.Validate(); err != nil {
			issues = append(issues, product.ID+": "+err.Error())
		}
		if previous, exists := seenSKUs[product.SKU]; exists {
			issues = append(issues, product.ID+": duplicate sku with "+previous)
		}
		seenSKUs[product.SKU] = product.ID
	}
	sort.Strings(issues)
	return issues, nil
}
