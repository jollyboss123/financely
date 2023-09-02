package currency

import (
	"github.com/google/uuid"
)

type Res struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	NumericCode string    `json:"numeric_code"`
	Fraction    int8      `json:"fraction"`
	Grapheme    string    `json:"grapheme"`
	Template    string    `json:"template"`
	Decimal     string    `json:"decimal"`
	Thousand    string    `json:"thousand"`
}

func Resource(currency *Schema) *Res {
	if currency == nil {
		return &Res{}
	}
	r := &Res{
		ID:          currency.ID,
		Code:        currency.Code,
		NumericCode: currency.NumericCode,
		Fraction:    currency.Fraction,
		Grapheme:    currency.Grapheme,
		Template:    currency.Template,
		Decimal:     currency.Decimal,
		Thousand:    currency.Thousand,
	}
	return r
}

func Resources(currencies []*Schema) []*Res {
	if len(currencies) == 0 {
		return make([]*Res, 0)
	}

	var r []*Res
	for _, c := range currencies {
		res := Resource(c)
		r = append(r, res)
	}
	return r
}
