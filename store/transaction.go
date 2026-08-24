package store

import (
	"errors"
	"strings"

	"go.etcd.io/bbolt"
)

type Transaction struct {
	tx       *bbolt.Tx
	writable bool
}

func (d *Database) Read(fn func(Transaction) error) error {
	return d.View(func(tx *bbolt.Tx) error { return fn(Transaction{tx: tx}) })
}

func (d *Database) Write(fn func(Transaction) error) error {
	return d.Update(func(tx *bbolt.Tx) error { return fn(Transaction{tx: tx, writable: true}) })
}

func (t Transaction) Put(bucket, key string, value any) error {
	if !t.writable {
		return errors.New("transaction is read-only")
	}
	return PutJSON(t.tx, bucket, key, value)
}

func (t Transaction) Get(bucket, key string, target any) error {
	return GetJSON(t.tx, bucket, key, target)
}

func (t Transaction) Remove(bucket, key string) error {
	if !t.writable {
		return errors.New("transaction is read-only")
	}
	return Delete(t.tx, bucket, key)
}

func (t Transaction) List(bucket string) ([]string, error) {
	return Keys(t.tx, bucket)
}

func (t Transaction) RequireBucket(bucket string) error {
	if strings.TrimSpace(bucket) == "" || !IsKnownBucket(bucket) {
		return errors.New("unknown bucket")
	}
	if t.tx.Bucket([]byte(bucket)) == nil {
		return errors.New("bucket missing")
	}
	return nil
}
