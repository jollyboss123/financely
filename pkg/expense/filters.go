package expense

import (
	"github.com/jollyboss123/finance-tracker/pkg/pagination"
	"net/url"
)

type Filter struct {
	Pagination      pagination.Filter
	Title           string `json:"title"`
	Amount          string `json:"amount"`
	TransactionDate string `json:"transaction_date"`
}

func Filters(queries url.Values) *Filter {
	p := pagination.New(queries)
	switch {
	case queries.Has("title"):
		fallthrough
	case queries.Has("transaction_date"):
		p.Search = true
	}

	return &Filter{
		Pagination:      *p,
		Title:           queries.Get("title"),
		Amount:          queries.Get("amount"),
		TransactionDate: queries.Get("transaction_date"),
	}
}
