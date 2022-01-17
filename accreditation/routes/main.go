package routes

import (
	"accreditation/app"
	"encoding/json"
	"net/http"
	"strings"
)

type Routes interface {
	Default() *http.ServeMux
}

type routes struct {
	accreditation app.Accreditation
	log           Logger
}

func healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}

func accounts(a app.Accreditation, log Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method == http.MethodPost {
			accountResponse, err := createAccountWithContext(ctx, r.Body, log, a)
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
		} else if r.Method == http.MethodGet {
			externalKey := strings.TrimPrefix(r.URL.Path, "/v1/accounts/")

			o, err := getAccountWithContext(ctx, externalKey, log, a)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if o == nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			res, err := json.Marshal(o)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write(res); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}

func (r *routes) Default() *http.ServeMux {
	middleware := http.NewServeMux()
	middleware.Handle("/v1/accounts/", accounts(r.accreditation, r.log))
	middleware.Handle("/v1/accounts", accounts(r.accreditation, r.log))
	middleware.Handle("/health", healthz())
	return middleware
}

func New(a app.Accreditation, log Logger) Routes {
	return &routes{
		accreditation: a,
		log:           log,
	}
}
