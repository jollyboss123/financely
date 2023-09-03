package expense

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jollyboss123/finance-tracker/pkg/currency"
	"github.com/jollyboss123/finance-tracker/pkg/rate"
	"github.com/jollyboss123/finance-tracker/pkg/server/message"
	"github.com/jollyboss123/finance-tracker/pkg/server/response"
	"github.com/jollyboss123/finance-tracker/pkg/validate"
	"github.com/shopspring/decimal"
	"net/http"
)

type Handler struct {
	expenseRepo  Expense
	currencyRepo currency.Currency
	exchangeRate *rate.ExchangeRates
	validator    *validator.Validate
}

func NewHandler(expenseRepo Expense, currencyRepo currency.Currency, exchangeRate *rate.ExchangeRates, validator *validator.Validate) *Handler {
	return &Handler{
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
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	cID, err := h.currencyRepo.ReadByCode(r.Context(), request.CurrencyCode)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	request.CurrencyID = cID
	request.BaseCurrencyCode = "MYR"
	request.BaseCurrencyID, err = h.currencyRepo.ReadByCode(r.Context(), request.BaseCurrencyCode)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	er, err := h.exchangeRate.GetRate(r.Context(), request.CurrencyCode, request.BaseCurrencyCode)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	amountDec := decimal.NewFromInt(request.Amount)
	rateDec := decimal.NewFromFloat(er)
	resultDec := amountDec.Mul(rateDec)
	request.BaseAmount = resultDec.IntPart()

	errs := validate.Validate(h.validator, request)
	if errs != nil {
		response.Errors(w, http.StatusBadRequest, errs)
		return
	}

	exID, err := h.expenseRepo.Create(r.Context(), &request)
	ex, err := h.expenseRepo.Read(r.Context(), exID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.Error(w, http.StatusBadRequest, message.ErrBadRequest)
			return
		}
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	e := Resource(ex)

	response.Json(w, http.StatusCreated, e)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	expenseID, err := uuid.Parse(chi.URLParam(r, "currencyID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, message.ErrBadRequest)
		return
	}

	ex, err := h.expenseRepo.Read(r.Context(), expenseID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.Error(w, http.StatusBadRequest, errors.New("no expense found for this ID"))
			return
		}
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	e := Resource(ex)

	response.Json(w, http.StatusOK, e)
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
				response.Error(w, http.StatusBadRequest, err)
				return
			}
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		expenses = resp

	default:
		resp, err := h.expenseRepo.List(ctx, filters)
		if err != nil {
			if errors.Is(err, ErrFetchingExpenses) {
				response.Error(w, http.StatusBadRequest, err)
				return
			}
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		expenses = resp
	}

	e := Resources(expenses)

	response.Json(w, http.StatusOK, e)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	expenseID, err := uuid.Parse(chi.URLParam(r, "expenseID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, message.ErrBadRequest)
		return
	}

	var request UpdateRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	request.ID = expenseID
	if request.CurrencyCode != "" {
		cID, err := h.currencyRepo.ReadByCode(r.Context(), request.CurrencyCode)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		request.CurrencyID = cID
	}
	request.BaseCurrencyCode = "MYR"
	request.BaseCurrencyID, err = h.currencyRepo.ReadByCode(r.Context(), request.BaseCurrencyCode)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	er, err := h.exchangeRate.GetRate(r.Context(), request.CurrencyCode, request.BaseCurrencyCode)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	amountDec := decimal.NewFromInt(request.Amount)
	rateDec := decimal.NewFromFloat(er)
	resultDec := amountDec.Mul(rateDec)
	request.BaseAmount = resultDec.IntPart()

	errs := validate.Validate(h.validator, request)
	if errs != nil {
		response.Errors(w, http.StatusBadRequest, errs)
		return
	}

	err = h.expenseRepo.Update(r.Context(), &request)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	ex, err := h.expenseRepo.Read(r.Context(), expenseID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	e := Resource(ex)

	response.Json(w, http.StatusOK, e)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	expenseID, err := uuid.Parse(chi.URLParam(r, "expenseID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	err = h.expenseRepo.Delete(r.Context(), expenseID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.Json(w, http.StatusOK, nil)
}

func (h *Handler) Total(w http.ResponseWriter, r *http.Request) {
	filters := Filters(r.URL.Query())
	var total int64

	total, err := h.expenseRepo.Total(r.Context(), filters)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Json(w, http.StatusOK, map[string]int64{
		"total": total,
	})
}
