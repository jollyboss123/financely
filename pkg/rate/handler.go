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
)

type Handler struct {
	rateRepo  Rate
	validator *validator.Validate
}

func NewHandler(rateRepo Rate, validator *validator.Validate) *Handler {
	return &Handler{
		rateRepo:  rateRepo,
		validator: validator,
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

	//cancel existing cron

	//run new cron
	for t := range cron.Cron(context.Background(), ur.startTime, ur.delay) {

		log.Println(t.Format("2006-01-02 15:04:05"))
	}
}
