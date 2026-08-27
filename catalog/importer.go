package catalog

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type ImportRow struct {
	Line        int
	ID          string
	SKU         string
	Name        string
	Category    string
	Unit        string
	PriceCents  int64
	Description string
}

type ImportResult struct {
	RowsAccepted int
	RowsRejected int
	Products     []Product
	Errors       []error
}

func ParseImport(reader io.Reader) (ImportResult, error) {
	if reader == nil {
		return ImportResult{}, errors.New("import reader is required")
	}
	scanner := bufio.NewScanner(reader)
	result := ImportResult{}
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if lineNumber == 1 && strings.EqualFold(strings.Split(line, ",")[0], "id") {
			continue
		}
		row, err := parseImportRow(line, lineNumber)
		if err != nil {
			result.RowsRejected++
			result.Errors = append(result.Errors, err)
			continue
		}
		product, err := NewProduct(row.ID, ProductDraft{SKU: row.SKU, Name: row.Name, Category: row.Category, Unit: row.Unit, PriceCents: row.PriceCents, Description: row.Description})
		if err != nil {
			result.RowsRejected++
			result.Errors = append(result.Errors, fmt.Errorf("line %d: %w", row.Line, err))
			continue
		}
		result.RowsAccepted++
		result.Products = append(result.Products, product)
	}
	if err := scanner.Err(); err != nil {
		return result, err
	}
	return result, nil
}

func parseImportRow(line string, number int) (ImportRow, error) {
	parts := strings.Split(line, ",")
	if len(parts) != 7 {
		return ImportRow{}, fmt.Errorf("line %d: expected seven columns", number)
	}
	price, err := strconv.ParseInt(strings.TrimSpace(parts[5]), 10, 64)
	if err != nil {
		return ImportRow{}, fmt.Errorf("line %d: invalid price", number)
	}
	return ImportRow{Line: number, ID: strings.TrimSpace(parts[0]), SKU: strings.TrimSpace(parts[1]), Name: strings.TrimSpace(parts[2]), Category: strings.TrimSpace(parts[3]), Unit: strings.TrimSpace(parts[4]), PriceCents: price, Description: strings.TrimSpace(parts[6])}, nil
}

func (s *Service) Import(reader io.Reader) (ImportResult, error) {
	result, err := ParseImport(reader)
	if err != nil {
		return result, err
	}
	for _, product := range result.Products {
		if saveErr := s.repo.Save(product); saveErr != nil {
			result.RowsRejected++
			result.RowsAccepted--
			result.Errors = append(result.Errors, fmt.Errorf("product %s: %w", product.ID, saveErr))
		}
	}
	return result, nil
}

func ExportCSV(products []Product) string {
	var builder strings.Builder
	builder.WriteString("id,sku,name,category,unit,price,description\n")
	for _, product := range SortProducts(products, "id", false) {
		builder.WriteString(strings.Join([]string{product.ID, product.SKU, product.Name, product.Category, product.Unit, strconv.FormatInt(product.PriceCents, 10), product.Description}, ","))
		builder.WriteByte('\n')
	}
	return builder.String()
}
