package reporting

import (
	"campgoods/catalog"
	"fmt"
	"sort"
	"strings"
)

func FormatSummary(summary Summary) string {
	status := summary.Availability
	lines := []string{
		fmt.Sprintf("%s [%s]", summary.Product.Name, summary.Product.SKU),
		fmt.Sprintf("category=%s status=%s stock=%d price=%s", summary.Product.Category, status, summary.Stock, summary.Product.DisplayPrice()),
		fmt.Sprintf("movements=%d price_changes=%d notes=%d", summary.MovementCount, summary.PriceChanges, len(summary.Notes)),
	}
	for _, note := range sortedNotes(summary.Notes) {
		marker := "-"
		if note.Pinned {
			marker = "*"
		}
		lines = append(lines, fmt.Sprintf("%s %s", marker, note.Text))
	}
	return strings.Join(lines, "\n")
}

func FormatDashboard(summaries []Summary, page catalog.ProductPage) string {
	lines := []string{fmt.Sprintf("page %d/%d total=%d", page.Page, page.TotalPages, page.Total)}
	for _, summary := range summaries {
		lines = append(lines, FormatSummary(summary))
	}
	return strings.Join(lines, "\n\n")
}

func sortedNotes(notes []ProductNote) []ProductNote {
	result := append([]ProductNote(nil), notes...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Pinned != result[j].Pinned {
			return result[i].Pinned
		}
		return result[i].ID < result[j].ID
	})
	return result
}
