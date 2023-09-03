package expense

import (
	"github.com/google/uuid"
	"time"
)

type Res struct {
	ID              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	Amount          int64     `json:"amount"`
	CurrencyCode    string    `json:"currency_code"`
	TransactionDate time.Time `json:"transaction_date"`
}

func Resource(expense *Schema) *Res {
	if expense == nil {
		return &Res{}
	}
	resource := &Res{
		ID:              expense.ID,
		Title:           expense.Title,
		Amount:          expense.Amount,
		CurrencyCode:    expense.CurrencyCode,
		TransactionDate: expense.TransactionDate,
	}

	return resource
}

func Resources(expenses []*Schema) []*Res {
	if len(expenses) == 0 {
		return make([]*Res, 0)
	}

	var resources []*Res
	for _, expense := range expenses {
		res := Resource(expense)
		resources = append(resources, res)
	}
	return resources
}
