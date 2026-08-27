package query

import (
	"errors"
	"strconv"
	"strings"

	"campgoods/catalog"
)

type ListRequest struct {
	Page       int
	PageSize   int
	Query      string
	Category   string
	Sort       string
	Descending bool
}

func (r ListRequest) Normalize() (ListRequest, error) {
	if r.Page == 0 {
		r.Page = 1
	}
	if r.PageSize == 0 {
		r.PageSize = 10
	}
	if r.Page < 1 || r.PageSize < 1 || r.PageSize > 100 {
		return ListRequest{}, errors.New("invalid pagination")
	}
	r.Query = strings.TrimSpace(r.Query)
	r.Category = strings.TrimSpace(r.Category)
	if r.Sort == "" {
		r.Sort = "id"
	}
	return r, nil
}

func (r ListRequest) Filter() catalog.ProductFilter {
	return catalog.ProductFilter{Query: r.Query, Category: r.Category, ActiveOnly: true}
}

func ParsePage(value string) (int, error) {
	if strings.TrimSpace(value) == "" {
		return 1, nil
	}
	page, err := strconv.Atoi(value)
	if err != nil || page < 1 {
		return 0, errors.New("page must be a positive integer")
	}
	return page, nil
}
