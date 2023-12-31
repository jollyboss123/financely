package health

import (
	"github.com/jollyboss123/finance-tracker/pkg/logger"
	"github.com/jollyboss123/finance-tracker/pkg/server/response"
	"net/http"
)

type Handler struct {
	logger     *logger.Logger
	healthRepo Repository
}

func NewHandler(logger *logger.Logger, health Repository) *Handler {
	return &Handler{
		logger:     logger,
		healthRepo: health,
	}
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Readiness(w http.ResponseWriter, _ *http.Request) {
	err := h.healthRepo.Readiness()
	if err != nil {
		response.Error(h.logger, w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
