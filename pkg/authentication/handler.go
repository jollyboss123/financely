package authentication

import (
	"errors"
	"github.com/alexedwards/argon2id"
	"github.com/go-playground/validator/v10"
	"github.com/jollyboss123/finance-tracker/pkg/logger"
	"github.com/jollyboss123/finance-tracker/pkg/middleware"
	"github.com/jollyboss123/finance-tracker/pkg/server/request"
	"github.com/jollyboss123/finance-tracker/pkg/server/response"
	"github.com/jollyboss123/finance-tracker/pkg/validate"
	"github.com/jollyboss123/scs/v2"
	"log"
	"net/http"
)

type Handler struct {
	logger    *logger.Logger
	validator *validator.Validate
	repo      User
	session   *scs.SessionManager
}

func NewHandler(l *logger.Logger, validator *validator.Validate, repo User, session *scs.SessionManager) *Handler {
	return &Handler{
		logger:    l,
		validator: validator,
		repo:      repo,
		session:   session,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		var mr *request.MalformedRequest
		if errors.As(err, &mr) {
			h.logger.Error().Err(err).Msg(mr.Msg)
			response.Error(h.logger, w, mr.Status, errors.New(mr.Msg))
			return
		}
		h.logger.Error().Err(err).Msg("failed decode register request")
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	errs := validate.Validate(h.validator, req)
	if errs != nil {
		h.logger.Error().Strs("error", errs).Msg("failed input validation")
		response.ValidationErrors(h.logger, w, errs)
		return
	}

	hashedPassword, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed hash password")
		response.Error(h.logger, w, http.StatusBadRequest, err)
		return
	}

	if err := h.repo.Register(r.Context(), req.FirstName, req.LastName, req.Email, hashedPassword); err != nil {
		h.logger.Error().Err(err).Msg("failed register user")
		response.Error(h.logger, w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		var mr *request.MalformedRequest
		if errors.As(err, &mr) {
			h.logger.Error().Err(err).Msg(mr.Msg)
			response.Error(h.logger, w, mr.Status, errors.New(mr.Msg))
			return
		}
		h.logger.Error().Err(err).Msg("failed decode login request")
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	errs := validate.Validate(h.validator, req)
	if errs != nil {
		h.logger.Error().Strs("error", errs).Msg("failed input validation")
		response.ValidationErrors(h.logger, w, errs)
		return
	}

	ctx := r.Context()

	user, _, err := h.repo.Login(ctx, &req)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed login")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := h.session.RenewToken(ctx); err != nil {
		h.logger.Error().Err(err).Msg("failed renew user token")
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	log.Println("putting userid now")
	//ctx = context.WithValue(ctx, middleware.KeyID, user.ID.String())
	h.session.Put(ctx, string(middleware.KeyID), user.ID.String())

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Protected(w http.ResponseWriter, _ *http.Request) {
	response.Json(h.logger, w, http.StatusOK, map[string]string{"success": "yup!"})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID := h.session.Get(r.Context(), string(middleware.KeyID))

	response.Json(h.logger, w, http.StatusOK, map[string]any{"user_id": userID})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	//currUser := h.session.Get(r.Context(), string(middleware.KeyID))

	err := h.session.Destroy(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("failed destroy session")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//userID, err := uuid.Parse(currUser.(string))
	//if err != nil {
	//	h.logger.Error().Err(err).Msg("failed parse userID")
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}

	//ok, err := h.repo.Logout(r.Context(), userID)
	//if err != nil {
	//	h.logger.Error().Err(err).Msg("failed logout")
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//
	//if !ok {
	//	response.Json(h.logger, w, http.StatusInternalServerError, map[string]string{"message": "unable to logout"})
	//}
}

func (h *Handler) Csrf(w http.ResponseWriter, r *http.Request) {
	userID := h.session.Get(r.Context(), string(middleware.KeyID))
	if userID == "" {
		h.logger.Error().Msg("user not logged in")
		response.Error(h.logger, w, http.StatusBadRequest, errors.New("you need to be logged in"))
		return
	}

	token, err := h.repo.Csrf(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("failed generate csrf token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Json(h.logger, w, http.StatusOK, map[string]string{"csrf_token": token})
}
