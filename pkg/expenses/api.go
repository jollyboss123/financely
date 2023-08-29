package expenses

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Expense struct {
	l *log.Logger
}

func NewExpense(l *log.Logger) *Expense {
	return &Expense{
		l: l,
	}
}

func (e *Expense) getExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hello world"))
	if err != nil {
		e.l.Fatalf("error writing message, %s", err)
	}
}

func (e *Expense) logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer e.l.Printf("request processed in %s\n", time.Now().Sub(startTime))
		next(w, r)
	}
}

func (e *Expense) SetupRoutes(mux *mux.Router) {
	mux.Methods(http.MethodGet).Subrouter()
	mux.HandleFunc("/", e.logger(e.getExpense))
}
