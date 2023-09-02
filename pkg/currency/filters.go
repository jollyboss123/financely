package currency

import (
	"github.com/jollyboss123/finance-tracker/pkg/pagination"
	"net/url"
)

type Filter struct {
	Pagination pagination.Filter
	Code       string `json:"code"`
}

func Filters(queries url.Values) *Filter {
	p := pagination.New(queries)
	if queries.Has("code") {
		p.Search = true
	}

	return &Filter{
		Pagination: *p,
		Code:       queries.Get("code"),
	}
}
