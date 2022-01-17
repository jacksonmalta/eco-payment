package server

import (
	"balance/routes"
	"fmt"
	"net/http"
	"time"
)

type Http struct {
	log Logger
	r   routes.Routes
}

func (h *Http) Start() {
	h.log.Info("Starting server")
	defaultRoutes := h.r.Default()
	server := &http.Server{
		Addr:         ":5003",
		Handler:      http.TimeoutHandler(defaultRoutes, 3*time.Second, "Timeout!!!"),
		ReadTimeout:  4 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	h.log.Info("Server is ready to handler request at :5003")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		h.log.Fatal(fmt.Sprintf("Could not listen on %s", err.Error()))
	}

}

func New(r routes.Routes, log Logger) *Http {
	return &Http{
		r:   r,
		log: log,
	}
}
