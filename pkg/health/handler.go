package health

import (
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
	healthRepo Repository
}

func NewHandler(health Repository) *Handler {
	return &Handler{
		healthRepo: health,
	}
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Readiness(w http.ResponseWriter, _ *http.Request) {
	err := h.healthRepo.Readiness()
	if err != nil {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusInternalServerError)

		p := map[string]string{
			"message": err.Error(),
		}
		data, err := json.Marshal(p)
		if err != nil {
			log.Println(err)
		}

		if string(data) == "null" {
			return
		}

		_, err = w.Write(data)
		if err != nil {
			log.Println(err)
		}
	}
}
