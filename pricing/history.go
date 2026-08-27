package pricing

import (
	"sort"

	"campgoods/store"
	"go.etcd.io/bbolt"
)

type History struct {
	db *store.Database
}

func NewHistory(db *store.Database) *History {
	return &History{db: db}
}

func (h *History) Append(change PriceChange) error {
	if err := change.Validate(); err != nil {
		return err
	}
	return h.db.Update(func(tx *bbolt.Tx) error { return store.PutJSON(tx, store.PricesBucket, change.ID, change) })
}

func (h *History) Find(id string) (PriceChange, error) {
	var change PriceChange
	err := h.db.View(func(tx *bbolt.Tx) error { return store.GetJSON(tx, store.PricesBucket, id, &change) })
	return change, err
}

func (h *History) ForProduct(productID string) ([]PriceChange, error) {
	var result []PriceChange
	err := h.db.View(func(tx *bbolt.Tx) error {
		keys, err := store.Keys(tx, store.PricesBucket)
		if err != nil {
			return err
		}
		for _, key := range keys {
			var change PriceChange
			if err := store.GetJSON(tx, store.PricesBucket, key, &change); err != nil {
				return err
			}
			if productID == "" || change.ProductID == productID {
				result = append(result, change)
			}
		}
		return nil
	})
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, err
}

func (h *History) All() ([]PriceChange, error) {
	return h.ForProduct("")
}
