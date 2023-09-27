package main

import (
	"net/http"

	ragDynamo "melkeydev/ragStackCDK/clients/dynamo"
	ragJWT "melkeydev/ragStackCDK/clients/jwt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type App struct {
	db  ragDynamo.UserStorageDB
	jwt ragJWT.TokenValidator
}

func NewApp(db ragDynamo.UserStorageDB, jwt ragJWT.TokenValidator) *App {
	return &App{
		db:  ragDynamo.NewDynamoDBClient(),
		jwt: ragJWT.NewJWTClient(db),
	}
}

func main() {
	db := ragDynamo.NewDynamoDBClient()
	jwt := ragJWT.NewJWTClient(db)

	app := NewApp(db, jwt)

	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/login":
			return app.LoginHandler(request)
		case "/register":
			return app.RegisterHandler(request)
		case "/refresh":
			return app.RefreshHandler(request)
		case "/test":
			return app.TestHandler(request)
		case "/protected":
			return ragJWT.ValidateJWTMiddleware(app.ProtectedHandler)(request)
		default:
			return events.APIGatewayProxyResponse{Body: "Not Found", StatusCode: http.StatusNotFound}, nil
		}
	})
}
