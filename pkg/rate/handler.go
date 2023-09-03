package rate

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/jollyboss123/finance-tracker/pkg/cron"
	"github.com/jollyboss123/finance-tracker/pkg/server/response"
	"github.com/jollyboss123/finance-tracker/pkg/validate"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	rateRepo  Rate
	validator *validator.Validate
	rates     *ExchangeRates
}

func NewHandler(rateRepo Rate, validator *validator.Validate, rates *ExchangeRates) *Handler {
	return &Handler{
		rateRepo:  rateRepo,
		validator: validator,
		rates:     rates,
	}
}

func (h *Handler) Reschedule(w http.ResponseWriter, r *http.Request) {
	var ur UpdateRequest
	err := json.NewDecoder(r.Body).Decode(&ur)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	errs := validate.Validate(h.validator, ur)
	if errs != nil {
		response.Errors(w, http.StatusBadRequest, errs)
		return
	}

	cron.Cancel("fetch.exchange-rates")

	jobFunc := func(t time.Time) {
		h.rates.GetRatesRemote(context.Background())
	}

	jobID := cron.Start("fetch.exchange-rates", ur.startTime, ur.delay, jobFunc)
	log.Printf("started new cron job: %s\n", jobID)
}
