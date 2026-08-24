package workflow

import (
	"errors"
	"fmt"
	"sort"

	"campgoods/catalog"
	"campgoods/reporting"
)

type BatchOutcome struct {
	Requested int
	Completed int
	Failures  []string
}

func (a *App) RegisterBatch(drafts map[string]catalog.ProductDraft) BatchOutcome {
	outcome := BatchOutcome{Requested: len(drafts)}
	ids := make([]string, 0, len(drafts))
	for id := range drafts {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		if _, err := a.RegisterProduct(id, drafts[id]); err != nil {
			outcome.Failures = append(outcome.Failures, id+": "+err.Error())
		} else {
			outcome.Completed++
		}
	}
	return outcome
}

func (a *App) ReceiveBatch(movements []BatchMovement) BatchOutcome {
	outcome := BatchOutcome{Requested: len(movements)}
	for _, movement := range movements {
		if _, err := a.ReceiveStock(movement.ID, movement.ProductID, movement.Quantity, movement.Reason); err != nil {
			outcome.Failures = append(outcome.Failures, movement.ID+": "+err.Error())
		} else {
			outcome.Completed++
		}
	}
	return outcome
}

type BatchMovement struct {
	ID, ProductID string
	Quantity      int64
	Reason        string
}

func (a *App) PinNote(noteID, productID, text string) error {
	if noteID == "" || productID == "" || text == "" {
		return errors.New("note fields are required")
	}
	return a.AddProductNote(reporting.ProductNote{ID: noteID, ProductID: productID, Text: text, Pinned: true})
}

func (a *App) WorkflowStatus(productID string) (string, error) {
	summary, err := a.ProductSummary(productID)
	if err != nil {
		return "", err
	}
	if !summary.Product.Active {
		return "archived", nil
	}
	if summary.Stock == 0 {
		return "replenishment-needed", nil
	}
	if summary.PriceChanges > 3 {
		return "price-review", nil
	}
	return "ready", nil
}

func (a *App) EnsureWorkflow(productID string) error {
	status, err := a.WorkflowStatus(productID)
	if err != nil {
		return err
	}
	if status == "archived" {
		return fmt.Errorf("%s is archived", productID)
	}
	return nil
}
