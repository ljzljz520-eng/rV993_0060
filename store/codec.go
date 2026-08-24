package store

import (
	"encoding/json"
	"errors"
	"sort"

	"go.etcd.io/bbolt"
)

func PutJSON(tx *bbolt.Tx, bucket, key string, value any) error {
	if tx == nil || !IsKnownBucket(bucket) || key == "" {
		return errors.New("invalid record target")
	}
	b := tx.Bucket([]byte(bucket))
	if b == nil {
		return errors.New("bucket missing")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return b.Put([]byte(key), data)
}

func GetJSON(tx *bbolt.Tx, bucket, key string, target any) error {
	if tx == nil || !IsKnownBucket(bucket) || key == "" {
		return errors.New("invalid record source")
	}
	b := tx.Bucket([]byte(bucket))
	if b == nil {
		return errors.New("bucket missing")
	}
	data := b.Get([]byte(key))
	if data == nil {
		return errors.New("record not found")
	}
	return json.Unmarshal(data, target)
}

func Delete(tx *bbolt.Tx, bucket, key string) error {
	if tx == nil || !IsKnownBucket(bucket) || key == "" {
		return errors.New("invalid record deletion")
	}
	return tx.Bucket([]byte(bucket)).Delete([]byte(key))
}

func Keys(tx *bbolt.Tx, bucket string) ([]string, error) {
	if tx == nil || !IsKnownBucket(bucket) {
		return nil, errors.New("invalid bucket")
	}
	var keys []string
	err := tx.Bucket([]byte(bucket)).ForEach(func(k, _ []byte) error {
		keys = append(keys, string(k))
		return nil
	})
	sort.Strings(keys)
	return keys, err
}

func Count(tx *bbolt.Tx, bucket string) (int, error) {
	keys, err := Keys(tx, bucket)
	return len(keys), err
}
