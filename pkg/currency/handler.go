package currency

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jollyboss123/finance-tracker/pkg/logger"
	"github.com/jollyboss123/finance-tracker/pkg/server/message"
	"github.com/jollyboss123/finance-tracker/pkg/server/response"
	"github.com/jollyboss123/finance-tracker/pkg/validate"
	"net/http"
)

type Handler struct {
	logger       *logger.Logger
	currencyRepo Currency
	validator    *validator.Validate
}

func NewHandler(logger *logger.Logger, currencyRepo Currency, validator *validator.Validate) *Handler {
	return &Handler{
		logger:       logger,
		currencyRepo: currencyRepo,
		validator:    validator,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var cr CreateRequest
	err := json.NewDecoder(r.Body).Decode(&cr)
	if err != nil {
		h.logger.Error().Err(err).Msgf("failed to decode request: %v", r.Body)
		response.Error(h.logger, w, http.StatusBadRequest, err)
		return
	}

	errs := validate.Validate(h.validator, cr)
	if errs != nil {
		h.logger.Error().Strs("error", errs).Msg("failed input validation")
		response.ValidationErrors(h.logger, w, errs)
		return
	}

	currID, err := h.currencyRepo.Create(r.Context(), &cr)
	c, err := h.currencyRepo.Read(r.Context(), currID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.logger.Error().Err(err).Msg("no rows found after currency created")
			response.Error(h.logger, w, http.StatusBadRequest, message.ErrBadRequest)
			return
		}
		h.logger.Error().Err(err).Msg("failed to read created currency")
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	h.logger.Info().Str("id", currID.String()).Msg("created currency")
	curr := Resource(c)

	response.Json(h.logger, w, http.StatusCreated, curr)
}

func (h *Handler) Fetch(w http.ResponseWriter, r *http.Request) {
	cID, err := uuid.Parse(chi.URLParam(r, "currencyID"))
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to decode currencyID")
		response.Error(h.logger, w, http.StatusBadRequest, message.ErrBadRequest)
		return
	}

	c, err := h.currencyRepo.Read(r.Context(), cID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.logger.Error().Str("id", cID.String()).Err(err).Msg("no currency found for this ID")
			response.Error(h.logger, w, http.StatusBadRequest, errors.New("no currency found for this ID"))
			return
		}
		h.logger.Error().Str("id", cID.String()).Err(err).Msg("failed to read currency")
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	curr := Resource(c)

	response.Json(h.logger, w, http.StatusOK, curr)
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
				h.logger.Error().Err(err).Msg("failed to fetch currency")
				response.Error(h.logger, w, http.StatusBadRequest, err)
				return
			}
			h.logger.Error().Err(err).Msg("failed to search currency")
			response.Error(h.logger, w, http.StatusInternalServerError, err)
			return
		}
		cs = resp
	default:
		resp, err := h.currencyRepo.List(ctx, f)
		if err != nil {
			if errors.Is(err, ErrFetchingCurrency) {
				h.logger.Error().Err(err).Msg("failed to fetch currency")
				response.Error(h.logger, w, http.StatusBadRequest, err)
				return
			}
			h.logger.Error().Err(err).Msg("failed to list currency")
			response.Error(h.logger, w, http.StatusInternalServerError, err)
			return
		}
		cs = resp
	}

	currs := Resources(cs)

	response.Json(h.logger, w, http.StatusOK, currs)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	cID, err := uuid.Parse(chi.URLParam(r, "currencyID"))
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to parse currencyID")
		response.Error(h.logger, w, http.StatusBadRequest, message.ErrBadRequest)
		return
	}

	var ur UpdateRequest
	err = json.NewDecoder(r.Body).Decode(&ur)
	if err != nil {
		h.logger.Error().Str("id", cID.String()).Err(err).Msgf("failed to decode: %v", r.Body)
		response.Error(h.logger, w, http.StatusBadRequest, err)
		return
	}
	ur.ID = cID

	errs := validate.Validate(h.validator, ur)
	if errs != nil {
		h.logger.Error().Str("id", ur.ID.String()).Strs("error", errs).Msg("failed input validation")
		response.ValidationErrors(h.logger, w, errs)
		return
	}

	err = h.currencyRepo.Update(r.Context(), &ur)
	if err != nil {
		h.logger.Error().Str("id", ur.ID.String()).Err(err).Msg("failed update currency")
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	c, err := h.currencyRepo.Read(r.Context(), cID)
	if err != nil {
		h.logger.Error().Str("id", ur.ID.String()).Err(err).Msg("failed read updated currency")
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	h.logger.Info().Str("id", ur.ID.String()).Msgf("updated currency: %v", c)
	curr := Resource(c)

	response.Json(h.logger, w, http.StatusOK, curr)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	cID, err := uuid.Parse(chi.URLParam(r, "currencyID"))
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to parse currencyID")
		response.Error(h.logger, w, http.StatusBadRequest, err)
		return
	}

	err = h.currencyRepo.Delete(r.Context(), cID)
	if err != nil {
		h.logger.Error().Str("id", cID.String()).Err(err).Msg("failed to delete currency")
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	h.logger.Info().Str("id", cID.String()).Msg("deleted currency")
	response.Json(h.logger, w, http.StatusOK, nil)
}
