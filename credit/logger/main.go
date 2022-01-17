package logger

import (
	"credit/app"
	"credit/authorizer"
	"credit/routes"
	"credit/server"
	"credit/settlement"
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

func New() (app.Logger, server.Logger, routes.Logger, authorizer.Logger, settlement.Logger) {
	return &logs{}, &logs{}, &logs{}, &logs{}, &logs{}
}
