package main

import (
	"credit/app"
	"credit/authorizer"
	"credit/logger"
	"credit/routes"
	"credit/server"
	"credit/services"
	"credit/settlement"
	"os"
)

func main() {
	logApp, logServer, logRoutes, logAuthorizer, logSettlement := logger.New()
	accreditationHttp, settlementHttp := services.NewHttp()
	confAuthorizer := &authorizer.Config{}
	confAuthorizer.WithUrl(os.Getenv("URL_ACCREDITATION"))
	accreditation := authorizer.New(logAuthorizer, confAuthorizer, accreditationHttp)
	confSettlement := &settlement.Config{}
	confSettlement.WithUrl(os.Getenv("URL_BALANCE"))
	balance := settlement.New(logSettlement, confSettlement, settlementHttp)
	credit := app.New(accreditation, balance, logApp)
	routes := routes.New(credit, logRoutes)
	serverHttp := server.New(routes, logServer)
	serverHttp.Start()
}
