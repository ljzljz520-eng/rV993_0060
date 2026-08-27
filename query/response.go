package query

import (
	"fmt"
	"strings"

	"campgoods/catalog"
)

type PageView struct {
	Page catalog.ProductPage
	Rows []string
}

func BuildPageView(page catalog.ProductPage) PageView {
	rows := make([]string, 0, len(page.Items))
	for index, product := range page.Items {
		rows = append(rows, fmt.Sprintf("%d. %s | %s | stock %d | %s", index+1, product.SKU, product.Name, product.Stock, product.DisplayPrice()))
	}
	return PageView{Page: page, Rows: rows}
}

func (v PageView) Render() string {
	lines := []string{fmt.Sprintf("Page %d of %d (%d products)", v.Page.Page, v.Page.TotalPages, v.Page.Total)}
	lines = append(lines, v.Rows...)
	if v.Page.HasPrevious {
		lines = append(lines, "previous available")
	}
	if v.Page.HasNext {
		lines = append(lines, "next available")
	}
	return strings.Join(lines, "\n")
}
