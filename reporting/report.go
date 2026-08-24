package reporting

import (
	"errors"
	"fmt"
	"strings"

	"campgoods/catalog"
	"campgoods/inventory"
	"campgoods/pricing"
	"campgoods/store"
	"go.etcd.io/bbolt"
)

type ProductNote struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Text      string `json:"text"`
	Pinned    bool   `json:"pinned"`
}

func (n ProductNote) Validate() error {
	if strings.TrimSpace(n.ID) == "" || strings.TrimSpace(n.ProductID) == "" {
		return errors.New("note identity is required")
	}
	if strings.TrimSpace(n.Text) == "" {
		return errors.New("note text is required")
	}
	return nil
}

type Summary struct {
	Product       catalog.Product
	Stock         int64
	MovementCount int
	PriceChanges  int
	Notes         []ProductNote
	Availability  string
}

type Service struct {
	products  *catalog.Service
	inventory *inventory.Service
	prices    *pricing.Service
	db        *store.Database
}

func NewService(products *catalog.Service, inventoryService *inventory.Service, prices *pricing.Service, db *store.Database) *Service {
	return &Service{products: products, inventory: inventoryService, prices: prices, db: db}
}

func (s *Service) AddNote(note ProductNote) error {
	if err := note.Validate(); err != nil {
		return err
	}
	if _, err := s.products.Get(note.ProductID); err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return store.PutJSON(tx, store.NotesBucket, note.ID, note) })
}

func (s *Service) Notes(productID string) ([]ProductNote, error) {
	var notes []ProductNote
	err := s.db.View(func(tx *bbolt.Tx) error {
		keys, err := store.Keys(tx, store.NotesBucket)
		if err != nil {
			return err
		}
		for _, key := range keys {
			var note ProductNote
			if err := store.GetJSON(tx, store.NotesBucket, key, &note); err != nil {
				return err
			}
			if productID == "" || note.ProductID == productID {
				notes = append(notes, note)
			}
		}
		return nil
	})
	return notes, err
}

func (s *Service) BuildSummary(productID string) (Summary, error) {
	product, err := s.products.Get(productID)
	if err != nil {
		return Summary{}, err
	}
	movements, err := s.inventory.History(productID)
	if err != nil {
		return Summary{}, err
	}
	changes, err := s.prices.HistoryFor(productID)
	if err != nil {
		return Summary{}, err
	}
	notes, err := s.Notes(productID)
	if err != nil {
		return Summary{}, err
	}
	availability := "out-of-stock"
	if !product.Active {
		availability = "archived"
	} else if product.Stock > 0 {
		availability = "available"
	}
	return Summary{Product: product, Stock: product.Stock, MovementCount: len(movements), PriceChanges: len(changes), Notes: notes, Availability: availability}, nil
}

func (s *Service) Dashboard(filter catalog.ProductFilter, page, size int) ([]Summary, catalog.ProductPage, error) {
	pageData, err := s.products.List(filter, page, size)
	if err != nil {
		return nil, catalog.ProductPage{}, err
	}
	result := make([]Summary, 0, len(pageData.Items))
	for _, product := range pageData.Items {
		summary, summaryErr := s.BuildSummary(product.ID)
		if summaryErr != nil {
			return nil, catalog.ProductPage{}, fmt.Errorf("summary %s: %w", product.ID, summaryErr)
		}
		result = append(result, summary)
	}
	return result, pageData, nil
}
