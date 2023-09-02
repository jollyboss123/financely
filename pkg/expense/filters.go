package expense

import (
	"github.com/jollyboss123/finance-tracker/pkg/pagination"
	"net/url"
)

type Filter struct {
	Pagination pagination.Filter
	Title      string `json:"title"`
	Amount     string `json:"amount"`
	Year       string `json:"year"`
	Month      string `json:"month"`
	Day        string `json:"day"`
	Base       string `json:"base"`
}

func Filters(queries url.Values) *Filter {
	p := pagination.New(queries)
	switch {
	case queries.Has("title"):
		fallthrough
	case queries.Has("year"):
		fallthrough
	case queries.Has("month"):
		fallthrough
	case queries.Has("day"):
		fallthrough
	case queries.Has("base"):
		p.Search = true
	}

	return &Filter{
		Pagination: *p,
		Title:      queries.Get("title"),
		Amount:     queries.Get("amount"),
		Year:       queries.Get("year"),
		Month:      queries.Get("month"),
		Day:        queries.Get("day"),
		Base:       queries.Get("base"),
	}
}
