package workflow

import (
	"errors"
	"fmt"

	"campgoods/catalog"
	"campgoods/query"
	"campgoods/reporting"
)

func (a *App) RegisterProduct(id string, draft catalog.ProductDraft) (catalog.Product, error) {
	if a == nil || a.Catalog == nil {
		return catalog.Product{}, errors.New("catalog unavailable")
	}
	return a.Catalog.Register(id, draft)
}

func (a *App) ReceiveStock(movementID, productID string, quantity int64, reason string) (string, error) {
	movement, err := a.Inventory.Record(movementID, productID, "receive", quantity, reason)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s stock=%d", movement.ProductID, movement.Balance), nil
}

func (a *App) ChangePrice(changeID, productID string, cents int64, reason string) (string, error) {
	change, err := a.Pricing.Change(changeID, productID, cents, reason, true)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s %s %d", change.ProductID, change.Direction(), change.After), nil
}

func (a *App) ListProducts(request query.ListRequest) (query.PageView, error) {
	normalized, err := request.Normalize()
	if err != nil {
		return query.PageView{}, err
	}
	page, err := a.Catalog.List(normalized.Filter(), normalized.Page, normalized.PageSize)
	if err != nil {
		return query.PageView{}, err
	}
	return query.BuildPageView(page), nil
}

func (a *App) AddProductNote(note reporting.ProductNote) error {
	return a.Reporting.AddNote(note)
}

func (a *App) ProductSummary(productID string) (reporting.Summary, error) {
	return a.Reporting.BuildSummary(productID)
}
