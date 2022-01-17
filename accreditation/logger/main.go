package logger

import (
	"accreditation/app"
	"accreditation/repository"
	"accreditation/routes"
	"accreditation/server"
	"log"
)

type logs struct{}

func (l *logs) Fatal(msg string) {
	log.Fatal(msg)
}

func (l *logs) Info(msg string) {
	log.Print(msg)
}

func (l *logs) Error(msg string) {
	log.Print(msg)
}

func New() (app.Logger, server.Logger, routes.Logger, repository.Logger) {
	return &logs{}, &logs{}, &logs{}, &logs{}
}
