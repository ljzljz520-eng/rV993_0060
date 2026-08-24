package catalog

import "strings"

type ProductFilter struct {
	Query       string
	Category    string
	ActiveOnly  bool
	InStockOnly bool
}

func (f ProductFilter) Normalize() ProductFilter {
	f.Query = strings.ToLower(strings.TrimSpace(f.Query))
	f.Category = strings.ToLower(strings.TrimSpace(f.Category))
	return f
}

func (f ProductFilter) Matches(p Product) bool {
	f = f.Normalize()
	if f.ActiveOnly && !p.Active {
		return false
	}
	if f.InStockOnly && p.Stock <= 0 {
		return false
	}
	if f.Category != "" && f.Category != p.Category {
		return false
	}
	if f.Query == "" {
		return true
	}
	return strings.Contains(strings.ToLower(p.Name), f.Query) || strings.Contains(strings.ToLower(p.SKU), f.Query) || strings.Contains(strings.ToLower(p.Description), f.Query)
}

func SortProducts(products []Product, field string, descending bool) []Product {
	result := append([]Product(nil), products...)
	for i := 1; i < len(result); i++ {
		current := result[i]
		j := i - 1
		for j >= 0 && comesAfter(result[j], current, field, descending) {
			result[j+1] = result[j]
			j--
		}
		result[j+1] = current
	}
	return result
}

func comesAfter(left, right Product, field string, descending bool) bool {
	var leftValue, rightValue string
	switch field {
	case "name":
		leftValue, rightValue = left.Name, right.Name
	case "category":
		leftValue, rightValue = left.Category, right.Category
	case "stock":
		leftValue, rightValue = formatNumber(left.Stock), formatNumber(right.Stock)
	default:
		leftValue, rightValue = left.ID, right.ID
	}
	if descending {
		return leftValue < rightValue
	}
	return leftValue > rightValue
}

func formatNumber(value int64) string {
	if value == 0 {
		return "0"
	}
	negative := value < 0
	if negative {
		value = -value
	}
	var digits []byte
	for value > 0 {
		digits = append(digits, byte('0'+value%10))
		value /= 10
	}
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	if negative {
		return "-" + string(digits)
	}
	return string(digits)
}
