package store

import (
	"errors"
	"path/filepath"
	"sync"

	"go.etcd.io/bbolt"
)

var ErrClosed = errors.New("store is closed")

type Database struct {
	mu   sync.RWMutex
	db   *bbolt.DB
	path string
}

func Open(path string) (*Database, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	db, err := bbolt.Open(filepath.Clean(path), 0600, nil)
	if err != nil {
		return nil, err
	}
	d := &Database{db: db, path: filepath.Clean(path)}
	if err := d.createBuckets(); err != nil {
		db.Close()
		return nil, err
	}
	return d, nil
}

func (d *Database) createBuckets() error {
	return d.db.Update(func(tx *bbolt.Tx) error {
		for _, name := range BucketNames() {
			if _, err := tx.CreateBucketIfNotExists([]byte(name)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (d *Database) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.db == nil {
		return ErrClosed
	}
	err := d.db.Close()
	d.db = nil
	return err
}

func (d *Database) Update(fn func(*bbolt.Tx) error) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.db == nil {
		return ErrClosed
	}
	return d.db.Update(fn)
}

func (d *Database) View(fn func(*bbolt.Tx) error) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if d.db == nil {
		return ErrClosed
	}
	return d.db.View(fn)
}

func (d *Database) Path() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.path
}

func (d *Database) IsOpen() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.db != nil
}
