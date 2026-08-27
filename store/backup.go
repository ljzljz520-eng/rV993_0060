package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.etcd.io/bbolt"
)

type BucketStats struct {
	Name    string
	Records int
}

func (d *Database) Stats() ([]BucketStats, error) {
	var stats []BucketStats
	err := d.View(func(tx *bbolt.Tx) error {
		for _, name := range BucketNames() {
			count, err := Count(tx, name)
			if err != nil {
				return err
			}
			stats = append(stats, BucketStats{Name: name, Records: count})
		}
		return nil
	})
	return stats, err
}

func (d *Database) CopyTo(path string) error {
	if path == "" {
		return errors.New("backup path is required")
	}
	if !d.IsOpen() {
		return ErrClosed
	}
	if filepath.Clean(path) == d.Path() {
		return errors.New("backup path must differ")
	}
	if err := os.MkdirAll(filepath.Dir(filepath.Clean(path)), 0755); err != nil {
		return err
	}
	return d.View(func(tx *bbolt.Tx) error {
		file, err := os.OpenFile(filepath.Clean(path), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := tx.WriteTo(file); err != nil {
			return fmt.Errorf("backup write: %w", err)
		}
		return file.Sync()
	})
}

func (d *Database) BucketStats(name string) (BucketStats, error) {
	if !IsKnownBucket(name) {
		return BucketStats{}, errors.New("unknown bucket")
	}
	var count int
	err := d.View(func(tx *bbolt.Tx) error { var err error; count, err = Count(tx, name); return err })
	return BucketStats{Name: name, Records: count}, err
}

func (d *Database) ClearBucket(name string) error {
	if !IsKnownBucket(name) {
		return errors.New("unknown bucket")
	}
	return d.Update(func(tx *bbolt.Tx) error {
		if err := tx.DeleteBucket([]byte(name)); err != nil && !errors.Is(err, bbolt.ErrBucketNotFound) {
			return err
		}
		_, err := tx.CreateBucket([]byte(name))
		return err
	})
}
