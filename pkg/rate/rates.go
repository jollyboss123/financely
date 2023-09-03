package rate

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
)

type ExchangeRates struct {
	rateRepo Rate
}

func NewExchangeRates(rateRepo Rate) *ExchangeRates {
	return &ExchangeRates{rateRepo: rateRepo}
}

func (er *ExchangeRates) ComputeRate(ctx context.Context, base, dest string) (float64, error) {
	br, err := er.rateRepo.Read(ctx, base)
	dr, err := er.rateRepo.Read(ctx, dest)
	if err != nil {
		return 0, err
	}
	return dr / br, nil
}

func (er *ExchangeRates) FetchRates(ctx context.Context) error {
	//TODO: change to call https://exchangeratesapi.io
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")

	if err != nil {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expect status code 200 got: %v", resp.Status)
	}
	defer resp.Body.Close()

	md := &cubes{}
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

	return nil
}

type cubes struct {
	CubeData []cube `xml:"Cube>Cube>Cube"`
}

type cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
