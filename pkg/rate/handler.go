package rate

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/jollyboss123/finance-tracker/config"
	"github.com/jollyboss123/finance-tracker/pkg/cron"
	"github.com/jollyboss123/finance-tracker/pkg/logger"
	"github.com/jollyboss123/finance-tracker/pkg/server/response"
	"github.com/jollyboss123/finance-tracker/pkg/validate"
	"net/http"
	"time"
)

type Handler struct {
	logger    *logger.Logger
	rateRepo  Rate
	validator *validator.Validate
	rates     *ExchangeRates
	cfg       *config.Config
}

func NewHandler(logger *logger.Logger, rateRepo Rate, validator *validator.Validate, rates *ExchangeRates, cfg *config.Config) *Handler {
	return &Handler{
		logger:    logger,
		rateRepo:  rateRepo,
		validator: validator,
		rates:     rates,
		cfg:       cfg,
	}
}

func (h *Handler) Reschedule(w http.ResponseWriter, r *http.Request) {
	var ur UpdateRequest
	err := json.NewDecoder(r.Body).Decode(&ur)
	if err != nil {
		h.logger.Error().Err(err).Msgf("failed to decode: %v", r.Body)
		response.Error(h.logger, w, http.StatusBadRequest, err)
		return
	}

	errs := validate.Validate(h.validator, ur)
	if errs != nil {
		h.logger.Error().Str("id", h.cfg.Cron.ExchangeRatesJobID).Strs("error", errs).Msg("failed input validation")
		response.ValidationErrors(h.logger, w, errs)
		return
	}

	cron.Cancel(h.cfg.Cron.ExchangeRatesJobID)
	h.logger.Info().Str("id", h.cfg.Cron.ExchangeRatesJobID).Msg("cancelled")

	jobFunc := func(t time.Time) {
		h.rates.GetRatesRemote(context.Background())
	}

	jobID, err := cron.Start(h.logger, h.cfg.Cron.ExchangeRatesJobID, ur.startTime, ur.delay, jobFunc)
	if err != nil {
		h.logger.Error().Str("id", h.cfg.Cron.ExchangeRatesJobID).Err(err).Msg("failed to start new job")
		response.Error(h.logger, w, http.StatusInternalServerError, err)
	}
	h.logger.Info().Str("id", jobID).Msgf("started new cron job at: %v with delay: %v", ur.startTime, ur.delay)
}
