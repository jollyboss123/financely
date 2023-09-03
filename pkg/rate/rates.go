package rate

import (
	"context"
	"encoding/xml"
	"fmt"
	s "github.com/shopspring/decimal"
	"net/http"
	"strconv"
)

type ExchangeRates struct {
	rateRepo Rate
}

func NewExchangeRates(rateRepo Rate) *ExchangeRates {
	return &ExchangeRates{rateRepo: rateRepo}
}

func (er *ExchangeRates) GetRate(ctx context.Context, base, dest string) (s.Decimal, error) {
	r, err := er.rateRepo.Read(ctx, base, dest)
	if err != nil {
		return s.Decimal{}, err
	}
	return r, nil
}

func (er *ExchangeRates) GetRatesRemote(ctx context.Context) error {
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")

	if err != nil {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expect status code 200 got: %v", resp.Status)
	}
	defer resp.Body.Close()

	md := &Cubes{}
	xml.NewDecoder(resp.Body).Decode(&md)
	for _, c := range md.CubeData {
		curr := c.Currency
		if curr != "USD" && curr != "JPY" && curr != "MYR" && curr != "SGD" &&
			curr != "AUD" && curr != "CAD" && curr != "CNY" && curr != "EUR" &&
			curr != "GBP" && curr != "HKD" && curr != "IDR" && curr != "KRW" &&
			curr != "TWD" && curr != "VND" {
			continue
		}

		r, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			fmt.Println("Error in ParseFloat:", err)
			return err
		}

		err = er.rateRepo.Create(ctx, "EUR", c.Currency, r)
		if err != nil {
			fmt.Println("Error in Create:", err)
			return err
		}
	}

	//e.rates["EUR"] = 1

	return nil
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
