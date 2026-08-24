package reporting

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"
)

type ExportBundle struct {
	GeneratedFor string    `json:"generated_for"`
	Summaries    []Summary `json:"summaries"`
}

func (s *Service) Export(productIDs []string) (ExportBundle, error) {
	if len(productIDs) == 0 {
		return ExportBundle{}, errors.New("at least one product is required")
	}
	ids := append([]string(nil), productIDs...)
	sort.Strings(ids)
	bundle := ExportBundle{GeneratedFor: strings.Join(ids, ","), Summaries: make([]Summary, 0, len(ids))}
	for _, id := range ids {
		summary, err := s.BuildSummary(id)
		if err != nil {
			return ExportBundle{}, err
		}
		bundle.Summaries = append(bundle.Summaries, summary)
	}
	return bundle, nil
}

func MarshalExport(bundle ExportBundle) ([]byte, error) {
	if len(bundle.Summaries) == 0 {
		return nil, errors.New("cannot marshal empty export")
	}
	return json.MarshalIndent(bundle, "", "  ")
}

func (s *Service) ExportActive() (ExportBundle, error) {
	products, err := s.products.ActiveProducts()
	if err != nil {
		return ExportBundle{}, err
	}
	ids := make([]string, 0, len(products))
	for _, product := range products {
		ids = append(ids, product.ID)
	}
	if len(ids) == 0 {
		return ExportBundle{}, errors.New("no active products")
	}
	return s.Export(ids)
}

func PlainText(bundle ExportBundle) string {
	lines := []string{bundle.GeneratedFor}
	for _, summary := range bundle.Summaries {
		lines = append(lines, FormatSummary(summary))
	}
	return strings.Join(lines, "\n\n")
}
