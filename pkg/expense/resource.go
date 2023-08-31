package expense

import s "github.com/shopspring/decimal"

type Res struct {
	ID     int       `json:"id"`
	Title  string    `json:"title"`
	Amount s.Decimal `json:"amount"`
}

func Resource(expense *Schema) *Res {
	if expense == nil {
		return &Res{}
	}
	resource := &Res{
		ID:     expense.ID,
		Title:  expense.Title,
		Amount: expense.Amount,
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
