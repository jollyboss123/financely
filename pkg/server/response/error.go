package response

import (
	"encoding/json"
	"github.com/jollyboss123/finance-tracker/pkg/logger"
	"net/http"
)

func Error(l *logger.Logger, w http.ResponseWriter, statusCode int, message error) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(statusCode)

	var p map[string]string
	if message == nil {
		write(l, w, nil)
		return
	}

	p = map[string]string{
		"message": message.Error(),
	}
	data, err := json.Marshal(p)
	if err != nil {
		l.Error().Err(err).Msg("failed to marshal error")
	}

	if string(data) == "null" {
		return
	}

	write(l, w, data)
}

func ValidationErrors(l *logger.Logger, w http.ResponseWriter, errors []string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(http.StatusUnprocessableEntity)

	if errors == nil {
		write(l, w, nil)
		return
	}

	p := map[string][]string{
		"message": errors,
	}
	data, err := json.Marshal(p)
	if err != nil {
		l.Error().Err(err).Msg("failed to marshal validation error")
	}

	if string(data) == "null" {
		return
	}

	write(l, w, data)
}

func write(l *logger.Logger, w http.ResponseWriter, data []byte) {
	_, err := w.Write(data)
	if err != nil {
		l.Error().Err(err).Msg("failed to write data")
	}
}
