package main

import (
	"debit/app"
	"debit/authorizer"
	"debit/logger"
	"debit/routes"
	"debit/server"
	"debit/services"
	"debit/settlement"
	"os"
)

func main() {
	logApp, logServer, logRoutes, logAuthorizer, logSettlement := logger.New()
	acdebitationHttp, settlementHttp := services.NewHttp()
	confAuthorizer := &authorizer.Config{}
	confAuthorizer.WithUrl(os.Getenv("URL_ACCREDITATION"))
	acdebitation := authorizer.New(logAuthorizer, confAuthorizer, acdebitationHttp)
	confSettlement := &settlement.Config{}
	confSettlement.WithUrl(os.Getenv("URL_BALANCE"))
	balance := settlement.New(logSettlement, confSettlement, settlementHttp)
	debit := app.New(acdebitation, balance, logApp)
	routes := routes.New(debit, logRoutes)
	serverHttp := server.New(routes, logServer)
	serverHttp.Start()
}
