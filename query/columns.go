package query

import (
	"errors"
	"strings"

	"campgoods/catalog"
)

type Column string

const (
	ColumnSKU      Column = "sku"
	ColumnName     Column = "name"
	ColumnCategory Column = "category"
	ColumnStock    Column = "stock"
	ColumnPrice    Column = "price"
)

func ParseColumns(value string) ([]Column, error) {
	if strings.TrimSpace(value) == "" {
		return []Column{ColumnSKU, ColumnName, ColumnCategory, ColumnStock, ColumnPrice}, nil
	}
	seen := make(map[Column]bool)
	columns := make([]Column, 0)
	for _, part := range strings.Split(value, ",") {
		column := Column(strings.TrimSpace(part))
		switch column {
		case ColumnSKU, ColumnName, ColumnCategory, ColumnStock, ColumnPrice:
		default:
			return nil, errors.New("unknown column")
		}
		if !seen[column] {
			seen[column] = true
			columns = append(columns, column)
		}
	}
	return columns, nil
}

func Value(product catalog.Product, column Column) string {
	switch column {
	case ColumnSKU:
		return product.SKU
	case ColumnName:
		return product.Name
	case ColumnCategory:
		return product.Category
	case ColumnStock:
		return formatInt(product.Stock)
	case ColumnPrice:
		return product.DisplayPrice()
	default:
		return ""
	}
}

func formatInt(value int64) string {
	if value == 0 {
		return "0"
	}
	negative := value < 0
	if negative {
		value = -value
	}
	result := make([]byte, 0, 20)
	for value > 0 {
		result = append(result, byte('0'+value%10))
		value /= 10
	}
	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}
	if negative {
		return "-" + string(result)
	}
	return string(result)
}
