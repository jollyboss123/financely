package expense

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jollyboss123/finance-tracker/pkg/server/message"
	"github.com/jollyboss123/finance-tracker/pkg/server/response"
	"github.com/jollyboss123/finance-tracker/pkg/validate"
	s "github.com/shopspring/decimal"
	"net/http"
	"strconv"
)

type Handler struct {
	expenseRepo Expense
	validator   *validator.Validate
}

func NewHandler(expenseRepo Expense, validator *validator.Validate) *Handler {
	return &Handler{
		expenseRepo: expenseRepo,
		validator:   validator,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var request CreateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

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
	expenseID, err := strconv.Atoi(chi.URLParam(r, "expenseID"))
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

	exs, err := h.expenseRepo.List(r.Context(), filters)
	if err != nil {
		if errors.Is(err, ErrFetchingExpenses) {
			response.Error(w, http.StatusBadRequest, err)
			return
		}
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	expenses := Resources(exs)

	response.Json(w, http.StatusOK, expenses)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	expenseID, err := strconv.Atoi(chi.URLParam(r, "expenseID"))
	var request UpdateRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
	}
	request.ID = expenseID

	errs := validate.Validate(h.validator, request)
	if errs != nil {
		response.Errors(w, http.StatusBadRequest, errs)
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
	expenseID, err := strconv.Atoi(chi.URLParam(r, "expenseID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
	}

	err = h.expenseRepo.Delete(r.Context(), expenseID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.Json(w, http.StatusOK, nil)
}

func (h *Handler) Total(w http.ResponseWriter, r *http.Request) {
	total, err := h.expenseRepo.Total(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Json(w, http.StatusOK, map[string]s.Decimal{
		"total": total,
	})
}
