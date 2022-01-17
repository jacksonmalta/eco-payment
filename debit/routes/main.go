package routes

import (
	"debit/app"
	"encoding/json"
	"net/http"
)

type Routes interface {
	Default() *http.ServeMux
}

type routes struct {
	debit app.Debit
	log   Logger
}

func healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}

func transactions(a app.Debit, log Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			ctx := r.Context()
			accountResponse, err := transactionWithContext(ctx, r.Body, log, a)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if accountResponse != nil {
				res, err := json.Marshal(accountResponse)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(accountResponse.Error.StatusCode)
				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write(res); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				return
			}

			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}

func (r *routes) Default() *http.ServeMux {
	middleware := http.NewServeMux()
	middleware.Handle("/v1/transactions", transactions(r.debit, r.log))
	middleware.Handle("/health", healthz())
	return middleware
}

func New(a app.Debit, log Logger) Routes {
	return &routes{
		debit: a,
		log:   log,
	}
}
