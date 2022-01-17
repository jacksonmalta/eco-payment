package main

import (
	"balance/app"
	"balance/logger"
	"balance/repository"
	"balance/routes"
	"balance/server"
	"balance/services"
	"os"
)

func main() {
	logApp, logServer, logRoutes, logDynamodb := logger.New()
	dynamodbService := services.NewDynamodb()
	dynamodbConfig := repository.Config{
		TableName: os.Getenv("TABLE_NAME"),
	}
	dynamodb := repository.NewDynamodb(dynamodbService, logDynamodb, dynamodbConfig)
	balance := app.New(dynamodb, logApp)
	routes := routes.New(balance, logRoutes)
	serverHttp := server.New(routes, logServer)
	serverHttp.Start()
}
