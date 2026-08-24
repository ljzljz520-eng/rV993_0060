package store

const (
	ProductsBucket  = "products"
	InventoryBucket = "inventory_movements"
	PricesBucket    = "price_changes"
	NotesBucket     = "product_notes"
)

func BucketNames() []string {
	return []string{ProductsBucket, InventoryBucket, PricesBucket, NotesBucket}
}

func BucketForEntity(entity string) string {
	switch entity {
	case "Product":
		return ProductsBucket
	case "InventoryMovement":
		return InventoryBucket
	case "PriceChange":
		return PricesBucket
	case "ProductNote":
		return NotesBucket
	default:
		return ""
	}
}

func IsKnownBucket(name string) bool {
	for _, candidate := range BucketNames() {
		if candidate == name {
			return true
		}
	}
	return false
}
