package inventory

import (
	"sort"

	"campgoods/store"
	"go.etcd.io/bbolt"
)

type Ledger struct {
	db *store.Database
}

func NewLedger(db *store.Database) *Ledger {
	return &Ledger{db: db}
}

func (l *Ledger) Append(movement InventoryMovement) error {
	if err := movement.Validate(); err != nil {
		return err
	}
	return l.db.Update(func(tx *bbolt.Tx) error { return store.PutJSON(tx, store.InventoryBucket, movement.ID, movement) })
}

func (l *Ledger) Find(id string) (InventoryMovement, error) {
	var movement InventoryMovement
	err := l.db.View(func(tx *bbolt.Tx) error { return store.GetJSON(tx, store.InventoryBucket, id, &movement) })
	return movement, err
}

func (l *Ledger) ForProduct(productID string) ([]InventoryMovement, error) {
	var result []InventoryMovement
	err := l.db.View(func(tx *bbolt.Tx) error {
		keys, err := store.Keys(tx, store.InventoryBucket)
		if err != nil {
			return err
		}
		for _, key := range keys {
			var movement InventoryMovement
			if err := store.GetJSON(tx, store.InventoryBucket, key, &movement); err != nil {
				return err
			}
			if productID == "" || movement.ProductID == productID {
				result = append(result, movement)
			}
		}
		return nil
	})
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, err
}

func (l *Ledger) All() ([]InventoryMovement, error) {
	return l.ForProduct("")
}
