package catalog

import (
	"errors"
	"sort"

	"campgoods/store"
	"go.etcd.io/bbolt"
)

type Repository struct {
	db *store.Database
}

func NewRepository(db *store.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(product Product) error {
	if err := product.Validate(); err != nil {
		return err
	}
	return r.db.Update(func(tx *bbolt.Tx) error { return store.PutJSON(tx, store.ProductsBucket, product.ID, product) })
}

func (r *Repository) Find(id string) (Product, error) {
	var product Product
	err := r.db.View(func(tx *bbolt.Tx) error { return store.GetJSON(tx, store.ProductsBucket, id, &product) })
	return product, err
}

func (r *Repository) Delete(id string) error {
	return r.db.Update(func(tx *bbolt.Tx) error { return store.Delete(tx, store.ProductsBucket, id) })
}

func (r *Repository) All() ([]Product, error) {
	var products []Product
	err := r.db.View(func(tx *bbolt.Tx) error {
		keys, err := store.Keys(tx, store.ProductsBucket)
		if err != nil {
			return err
		}
		products = make([]Product, 0, len(keys))
		for _, key := range keys {
			var product Product
			if err := store.GetJSON(tx, store.ProductsBucket, key, &product); err != nil {
				return err
			}
			products = append(products, product)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(products, func(i, j int) bool { return products[i].ID < products[j].ID })
	return products, nil
}

func (r *Repository) Require(id string) (Product, error) {
	product, err := r.Find(id)
	if err != nil {
		return Product{}, err
	}
	if product.ID == "" {
		return Product{}, errors.New("product has no identity")
	}
	return product, nil
}
