package currency

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jollyboss123/finance-tracker/pkg/server/message"
	"github.com/jollyboss123/finance-tracker/pkg/server/response"
	"github.com/jollyboss123/finance-tracker/pkg/validate"
	"net/http"
)

type Handler struct {
	currencyRepo Currency
	validator    *validator.Validate
}

func NewHandler(currencyRepo Currency, validator *validator.Validate) *Handler {
	return &Handler{
		currencyRepo: currencyRepo,
		validator:    validator,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var cr CreateRequest
	err := json.NewDecoder(r.Body).Decode(&cr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	errs := validate.Validate(h.validator, cr)
	if errs != nil {
		response.Errors(w, http.StatusBadRequest, errs)
		return
	}

	currID, err := h.currencyRepo.Create(r.Context(), &cr)
	c, err := h.currencyRepo.Read(r.Context(), currID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.Error(w, http.StatusBadRequest, message.ErrBadRequest)
			return
		}
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	curr := Resource(c)

	response.Json(w, http.StatusCreated, curr)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	cID, err := uuid.Parse(chi.URLParam(r, "currencyID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, message.ErrBadRequest)
		return
	}

	c, err := h.currencyRepo.Read(r.Context(), cID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.Error(w, http.StatusBadRequest, errors.New("no currency found for this ID"))
			return
		}
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	curr := Resource(c)

	response.Json(w, http.StatusOK, curr)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	f := Filters(r.URL.Query())

	var cs []*Schema
	ctx := r.Context()

	switch f.Pagination.Search {
	case true:
		resp, err := h.currencyRepo.Search(ctx, f)
		if err != nil {
			if errors.Is(err, ErrFetchingCurrency) {
				response.Error(w, http.StatusBadRequest, err)
				return
			}
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		cs = resp
	default:
		resp, err := h.currencyRepo.List(ctx, f)
		if err != nil {
			if errors.Is(err, ErrFetchingCurrency) {
				response.Error(w, http.StatusBadRequest, err)
				return
			}
			response.Error(w, http.StatusInternalServerError, err)
			return
		}
		cs = resp
	}

	currs := Resources(cs)

	response.Json(w, http.StatusOK, currs)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	cID, err := uuid.Parse(chi.URLParam(r, "currencyID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, message.ErrBadRequest)
		return
	}

	var ur UpdateRequest
	err = json.NewDecoder(r.Body).Decode(&ur)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	ur.ID = cID

	errs := validate.Validate(h.validator, ur)
	if errs != nil {
		response.Errors(w, http.StatusBadRequest, errs)
		return
	}

	err = h.currencyRepo.Update(r.Context(), &ur)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	c, err := h.currencyRepo.Read(r.Context(), cID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	curr := Resource(c)

	response.Json(w, http.StatusOK, curr)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	cID, err := uuid.Parse(chi.URLParam(r, "currencyID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}

	err = h.currencyRepo.Delete(r.Context(), cID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	response.Json(w, http.StatusOK, nil)
}
