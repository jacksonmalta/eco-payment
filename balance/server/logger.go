package server

type Logger interface {
	Info(msg string)
	Error(msg string)
	Fatal(msg string)
}
