package workflow

import (
	"errors"
	"path/filepath"

	"campgoods/catalog"
	"campgoods/inventory"
	"campgoods/pricing"
	"campgoods/reporting"
	"campgoods/store"
)

type App struct {
	DB        *store.Database
	Catalog   *catalog.Service
	Inventory *inventory.Service
	Pricing   *pricing.Service
	Reporting *reporting.Service
}

func Open(path string) (*App, error) {
	if filepath.Clean(path) == "." || path == "" {
		return nil, errors.New("storage path is required")
	}
	db, err := store.Open(path)
	if err != nil {
		return nil, err
	}
	products := catalog.NewService(catalog.NewRepository(db))
	ledger := inventory.NewLedger(db)
	history := pricing.NewHistory(db)
	inventoryService := inventory.NewService(products, ledger)
	pricingService := pricing.NewService(products, history)
	return &App{DB: db, Catalog: products, Inventory: inventoryService, Pricing: pricingService, Reporting: reporting.NewService(products, inventoryService, pricingService, db)}, nil
}

func NewWithDatabase(db *store.Database) *App {
	products := catalog.NewService(catalog.NewRepository(db))
	ledger := inventory.NewLedger(db)
	history := pricing.NewHistory(db)
	inventoryService := inventory.NewService(products, ledger)
	pricingService := pricing.NewService(products, history)
	return &App{DB: db, Catalog: products, Inventory: inventoryService, Pricing: pricingService, Reporting: reporting.NewService(products, inventoryService, pricingService, db)}
}

func (a *App) Close() error {
	if a == nil || a.DB == nil {
		return errors.New("app is not open")
	}
	return a.DB.Close()
}
