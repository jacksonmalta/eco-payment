package main

import (
	"accreditation/app"
	"accreditation/logger"
	"accreditation/repository"
	"accreditation/routes"
	"accreditation/server"
	"accreditation/services"
	"os"
)

func main() {
	logApp, logServer, logRoutes, logDynamodb := logger.New()
	dynamodbService := services.NewDynamodb()
	dynamodbConfig := repository.Config{
		TableName: os.Getenv("TABLE_NAME"),
	}
	dynamodb := repository.NewDynamodb(dynamodbService, logDynamodb, dynamodbConfig)
	accreditation := app.New(dynamodb, logApp)
	routes := routes.New(accreditation, logRoutes)
	serverHttp := server.New(routes, logServer)
	serverHttp.Start()
}
