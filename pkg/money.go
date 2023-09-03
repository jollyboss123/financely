package pkg

import "github.com/jollyboss123/finance-tracker/pkg/currency"

type Money struct {
	amount   int64
	currency *currency.Schema
}

func (m *Money) Display(f currency.Formatter, amount int64) string {
	return f.Format(amount)
}
