package store

import (
	"path/filepath"
	"testing"
)

func TestDatabaseStatsAndTransaction(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "store.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.Write(func(tx Transaction) error { return tx.Put(ProductsBucket, "id", map[string]string{"id": "id"}) }); err != nil {
		t.Fatal(err)
	}
	var value map[string]string
	if err := db.Read(func(tx Transaction) error { return tx.Get(ProductsBucket, "id", &value) }); err != nil {
		t.Fatal(err)
	}
	if value["id"] != "id" {
		t.Fatalf("value=%v", value)
	}
	stats, err := db.Stats()
	if err != nil || len(stats) != 4 {
		t.Fatalf("stats=%v err=%v", stats, err)
	}
}
