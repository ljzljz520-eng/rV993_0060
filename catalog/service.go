package catalog

import (
	"errors"
	"fmt"
)

type ProductPage struct {
	Items       []Product
	Page        int
	PageSize    int
	Total       int
	TotalPages  int
	HasPrevious bool
	HasNext     bool
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(id string, draft ProductDraft) (Product, error) {
	if _, err := s.repo.Find(id); err == nil {
		return Product{}, fmt.Errorf("product %s already exists", id)
	}
	product, err := NewProduct(id, draft)
	if err != nil {
		return Product{}, err
	}
	if err := s.repo.Save(product); err != nil {
		return Product{}, err
	}
	return product, nil
}

func (s *Service) Get(id string) (Product, error) {
	return s.repo.Require(id)
}

func (s *Service) UpdateDescription(id, description string) (Product, error) {
	product, err := s.repo.Require(id)
	if err != nil {
		return Product{}, err
	}
	if description == "" {
		return Product{}, errors.New("description cannot be empty")
	}
	product.Description = description
	product.Version++
	if err := s.repo.Save(product); err != nil {
		return Product{}, err
	}
	return product, nil
}

func (s *Service) SetActive(id string, active bool) (Product, error) {
	product, err := s.repo.Require(id)
	if err != nil {
		return Product{}, err
	}
	product.Active = active
	product.Version++
	if err := s.repo.Save(product); err != nil {
		return Product{}, err
	}
	return product, nil
}

func (s *Service) List(filter ProductFilter, page, pageSize int) (ProductPage, error) {
	if page < 1 || pageSize < 1 {
		return ProductPage{}, errors.New("page and page size must be positive")
	}
	products, err := s.repo.All()
	if err != nil {
		return ProductPage{}, err
	}
	filtered := make([]Product, 0, len(products))
	for _, product := range products {
		if filter.Matches(product) {
			filtered = append(filtered, product)
		}
	}
	filtered = SortProducts(filtered, "id", false)
	totalPages := (len(filtered) + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}
	start := (page - 1) * pageSize
	if page > 1 {
		start += pageSize
	}
	if start > len(filtered) {
		start = len(filtered)
	}
	end := start + pageSize
	if end > len(filtered) {
		end = len(filtered)
	}
	items := append([]Product(nil), filtered[start:end]...)
	return ProductPage{Items: items, Page: page, PageSize: pageSize, Total: len(filtered), TotalPages: totalPages, HasPrevious: page > 1, HasNext: page < totalPages}, nil
}

func (s *Service) Archive(id string) (Product, error) {
	return s.SetActive(id, false)
}

func (s *Service) Restore(id string) (Product, error) {
	return s.SetActive(id, true)
}

func (s *Service) SaveForInventory(product Product) error {
	return s.repo.Save(product)
}
