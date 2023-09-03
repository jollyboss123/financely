package expense

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jollyboss123/finance-tracker/pkg/currency"
	"github.com/jollyboss123/finance-tracker/pkg/logger"
	"github.com/jollyboss123/finance-tracker/pkg/rate"
	"github.com/jollyboss123/finance-tracker/pkg/server/message"
	"github.com/jollyboss123/finance-tracker/pkg/server/response"
	"github.com/jollyboss123/finance-tracker/pkg/validate"
	"github.com/shopspring/decimal"
	"net/http"
	"strings"
)

type Handler struct {
	logger       *logger.Logger
	expenseRepo  Expense
	currencyRepo currency.Currency
	exchangeRate *rate.ExchangeRates
	validator    *validator.Validate
}

func NewHandler(
	logger *logger.Logger,
	expenseRepo Expense,
	currencyRepo currency.Currency,
	exchangeRate *rate.ExchangeRates,
	validator *validator.Validate,
) *Handler {
	return &Handler{
		logger:       logger,
		expenseRepo:  expenseRepo,
		currencyRepo: currencyRepo,
		exchangeRate: exchangeRate,
		validator:    validator,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var request CreateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to decode request")
		response.Error(h.logger, w, http.StatusBadRequest, err)
		return
	}

	cID, err := h.currencyRepo.ReadByCode(r.Context(), request.CurrencyCode)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get currencyID")
		response.Error(h.logger, w, http.StatusBadRequest, err)
		return
	}
	request.CurrencyID = cID
	//TODO: get user base currency
	request.BaseCurrencyCode = "MYR"
	request.BaseCurrencyID, err = h.currencyRepo.ReadByCode(r.Context(), request.BaseCurrencyCode)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get base currencyID")
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	er, err := h.exchangeRate.ComputeRate(r.Context(), request.CurrencyCode, request.BaseCurrencyCode)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get exchange rate")
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}
	amountDec := decimal.NewFromInt(request.Amount)
	rateDec := decimal.NewFromFloat(er)
	resultDec := amountDec.Mul(rateDec)
	request.BaseAmount = resultDec.IntPart()

	errs := validate.Validate(h.validator, request)
	if errs != nil {
		h.logger.Error().Strs("error", errs).Msg("failed input validation")
		response.ValidationErrors(h.logger, w, errs)
		return
	}

	exID, err := h.expenseRepo.Create(r.Context(), &request)
	ex, err := h.expenseRepo.Read(r.Context(), exID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.logger.Error().Err(err).Msg("no rows found after expense created")
			response.Error(h.logger, w, http.StatusBadRequest, message.ErrBadRequest)
			return
		}
		h.logger.Error().Err(err).Msg("failed read created expense")
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	h.logger.Info().Str("id", ex.ID.String()).Msg("new expense created")
	e := Resource(ex)

	response.Json(h.logger, w, http.StatusCreated, e)
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	expenseID, err := uuid.Parse(chi.URLParam(r, "currencyID"))
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to parse currencyID")
		response.Error(h.logger, w, http.StatusBadRequest, message.ErrBadRequest)
		return
	}

	ex, err := h.expenseRepo.Read(r.Context(), expenseID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.logger.Error().Err(err).Msgf("no expense found for this currencyID: %v", expenseID)
			response.Error(h.logger, w, http.StatusBadRequest, errors.New("no expense found for this ID"))
			return
		}
		h.logger.Error().Err(err).Msgf("failed read expense for this currencyID: %v", expenseID)
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	e := Resource(ex)

	response.Json(h.logger, w, http.StatusOK, e)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	filters := Filters(r.URL.Query())

	var expenses []*Schema
	ctx := r.Context()

	switch filters.Pagination.Search {
	case true:
		resp, err := h.expenseRepo.Search(ctx, filters)
		if err != nil {
			if errors.Is(err, ErrFetchingExpenses) {
				h.logger.Error().Str("search", "true").Err(err).Msg("failed fetch expenses")
				response.Error(h.logger, w, http.StatusBadRequest, err)
				return
			}
			h.logger.Error().Str("search", "true").Err(err).Msg("failed to search expenses")
			response.Error(h.logger, w, http.StatusInternalServerError, err)
			return
		}
		expenses = resp

	default:
		resp, err := h.expenseRepo.List(ctx, filters)
		if err != nil {
			if errors.Is(err, ErrFetchingExpenses) {
				h.logger.Error().Str("search", "false").Err(err).Msg("failed fetch expenses")
				response.Error(h.logger, w, http.StatusBadRequest, err)
				return
			}
			h.logger.Error().Str("search", "false").Err(err).Msg("failed to list expenses")
			response.Error(h.logger, w, http.StatusInternalServerError, err)
			return
		}
		expenses = resp
	}

	e := Resources(expenses)

	response.Json(h.logger, w, http.StatusOK, e)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	expenseID, err := uuid.Parse(chi.URLParam(r, "expenseID"))
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to parse expenseID")
		response.Error(h.logger, w, http.StatusBadRequest, message.ErrBadRequest)
		return
	}

	var request UpdateRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error().Str("id", expenseID.String()).Err(err).Msg("failed to decode request")
		response.Error(h.logger, w, http.StatusBadRequest, err)
		return
	}
	request.ID = expenseID
	if request.CurrencyCode != "" {
		cID, err := h.currencyRepo.ReadByCode(r.Context(), request.CurrencyCode)
		if err != nil {
			h.logger.Error().Str("id", request.ID.String()).Err(err).Msgf("failed to fetch currency by code: %s", request.CurrencyCode)
			response.Error(h.logger, w, http.StatusBadRequest, err)
			return
		}
		request.CurrencyID = cID
	}
	//TODO: get user base currency
	request.BaseCurrencyCode = "MYR"
	request.BaseCurrencyID, err = h.currencyRepo.ReadByCode(r.Context(), request.BaseCurrencyCode)
	if err != nil {
		h.logger.Error().Str("id", request.ID.String()).Err(err).Msgf("failed to fetch base currency by code: %s", request.BaseCurrencyCode)
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	er, err := h.exchangeRate.ComputeRate(r.Context(), request.CurrencyCode, request.BaseCurrencyCode)
	if err != nil {
		h.logger.Error().Str("id", request.ID.String()).Err(err).Msg("failed to get exchange rate")
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}
	amountDec := decimal.NewFromInt(request.Amount)
	rateDec := decimal.NewFromFloat(er)
	resultDec := amountDec.Mul(rateDec)
	request.BaseAmount = resultDec.IntPart()

	errs := validate.Validate(h.validator, request)
	if errs != nil {
		h.logger.Error().Str("id", request.ID.String()).Strs("error", errs).Msg("failed input validation")
		response.ValidationErrors(h.logger, w, errs)
		return
	}

	err = h.expenseRepo.Update(r.Context(), &request)
	if err != nil {
		h.logger.Error().Str("id", request.ID.String()).Err(err).Msg("failed to update expense")
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}
	ex, err := h.expenseRepo.Read(r.Context(), expenseID)
	if err != nil {
		h.logger.Error().Str("id", request.ID.String()).Err(err).Msg("failed to read updated expense")
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	h.logger.Info().Str("id", request.ID.String()).Msgf("updated expense to: %v", ex)
	e := Resource(ex)

	response.Json(h.logger, w, http.StatusOK, e)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	expenseID, err := uuid.Parse(chi.URLParam(r, "expenseID"))
	if err != nil {
		h.logger.Error().Str("id", expenseID.String()).Err(err).Msg("failed to parse expenseID")
		response.Error(h.logger, w, http.StatusBadRequest, err)
		return
	}

	err = h.expenseRepo.Delete(r.Context(), expenseID)
	if err != nil {
		h.logger.Error().Str("id", expenseID.String()).Err(err).Msg("failed to delete expense")
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	h.logger.Info().Str("id", expenseID.String()).Msg("deleted expense")
	response.Json(h.logger, w, http.StatusOK, nil)
}

func (h *Handler) Total(w http.ResponseWriter, r *http.Request) {
	filters := Filters(r.URL.Query())
	if filters.Currency != "" {
		filters.Currency = strings.ToUpper(filters.Currency)
	}
	var total int64

	total, err := h.expenseRepo.Total(r.Context(), filters)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get total expense")
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	//TODO: get user base currency
	baseCC := "MYR"

	er, err := h.exchangeRate.ComputeRate(r.Context(), baseCC, filters.Currency)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get exchange rate")
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	amountDec := decimal.NewFromInt(total)
	rateDec := decimal.NewFromFloat(er)
	resultDec := amountDec.Mul(rateDec)
	total = resultDec.IntPart()

	response.Json(h.logger, w, http.StatusOK, map[string]int64{
		"total": total,
	})
}

func (h *Handler) Average(w http.ResponseWriter, r *http.Request) {
	filters := Filters(r.URL.Query())
	if filters.Currency != "" {
		filters.Currency = strings.ToUpper(filters.Currency)
	}
	var avg int64

	avg, err := h.expenseRepo.Average(r.Context(), filters)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get average expense")
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	//TODO: get user base currency
	baseCC := "MYR"

	er, err := h.exchangeRate.ComputeRate(r.Context(), baseCC, filters.Currency)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get exchange rate")
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	amountDec := decimal.NewFromInt(avg)
	rateDec := decimal.NewFromFloat(er)
	resultDec := amountDec.Mul(rateDec)
	avg = resultDec.IntPart()

	response.Json(h.logger, w, http.StatusOK, map[string]int64{
		"average": avg,
	})
}
