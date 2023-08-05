package main

import (
	"encoding/json"
	"log"
	"net/http"

	ragCrypto "melkeydev/ragStackCDK/clients/crypto"
	ragDynamo "melkeydev/ragStackCDK/clients/dynamo"
	ragJWT "melkeydev/ragStackCDK/clients/jwt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"name"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	userTableName = "user-table-name"
)

func RegisterHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var registerReq RegisterRequest

	err := json.Unmarshal([]byte(request.Body), &registerReq)
	if err != nil {
		// I want to use a different loggign package
		log.Printf("Unable to unmarshal register request:%v", err)
		return events.APIGatewayProxyResponse{Body: "Invalid Request"}, err
	}

	// Validate if the user already exists in DynamoDB
	// TODO: This should be something like email and/or username
	_, err = ragDynamo.GetUserFromDynamoDB(registerReq.Username)
	if err == nil {
		return events.APIGatewayProxyResponse{Body: "Username already exists", StatusCode: http.StatusBadRequest}, nil
	}

	hashedPassword, err := ragCrypto.GeneratePassword(registerReq.Password)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Invalid request", StatusCode: http.StatusBadRequest}, nil
	}

	// token, err := generateToken(registerReq.Username)
	token, err := ragJWT.GenerateToken(registerReq.Username)
	if err != nil {
		log.Print("Could not issue jwt token")
		return events.APIGatewayProxyResponse{Body: "Internal Server Error - Generating token", StatusCode: 500}, err
	}

	// Store the username and hashed password in DynamoDB
	err = ragDynamo.AddUserToDynamoDB(registerReq.Username, string(hashedPassword), token)
	if err != nil {
		log.Printf("Failed to add user to DynamoDB: %v", err)
		return events.APIGatewayProxyResponse{Body: "Internal Server Error - DDB", StatusCode: http.StatusInternalServerError}, nil
	}

	responseBody := map[string]string{
		"token": token,
	}

	responsejson, err := json.Marshal(responseBody)

	if err != nil {
		log.Printf("Failed to marshal response %v", err)
	}
	return events.APIGatewayProxyResponse{Body: string(responsejson), StatusCode: 200}, nil
}

func LoginHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var loginReq LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginReq)
	if err != nil {
		log.Printf("Unable to unmarshal login request:%v", err)
		return events.APIGatewayProxyResponse{Body: "Invalid Request"}, err
	}

	// Need server side validaton from inputs
	user, err := ragDynamo.GetUserFromDynamoDB(loginReq.Username)
	if err != nil {
		log.Printf("Failed to retrieve user from DynamoDB: %v", err)
		return events.APIGatewayProxyResponse{Body: "Invalid username or password", StatusCode: http.StatusUnauthorized}, nil
	}

	// Validate password
	if !ragCrypto.ComparePasswords(user.Password, loginReq.Password) {
		return events.APIGatewayProxyResponse{Body: "Invalid username or password", StatusCode: http.StatusUnauthorized}, nil
	}

	// Passwords match, generate and store a new JWT token
	token, err := ragJWT.GenerateToken(loginReq.Password)
	if err != nil {
		log.Print("Could not issue jwt token")
		return events.APIGatewayProxyResponse{Body: "Internal Server Error - Generating token", StatusCode: 500}, err
	}

	// Update the user's token in the database
	err = ragDynamo.UpdateUserTokenInDynamoDB(loginReq.Username, token)
	if err != nil {
		log.Printf("Failed to update user's token in DynamoDB: %v", err)
		return events.APIGatewayProxyResponse{Body: "Internal Server Error", StatusCode: http.StatusInternalServerError}, nil
	}

	responseBody := map[string]string{
		"token": token,
	}

	responsejson, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("Failed to marshal response %v", err)
	}
	return events.APIGatewayProxyResponse{Body: string(responsejson), StatusCode: 200}, nil

}

func main() {
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/login":
			return LoginHandler(request)
		case "/register":
			return RegisterHandler(request)
		default:
			return events.APIGatewayProxyResponse{Body: "Not Found", StatusCode: http.StatusNotFound}, nil
		}
	})
}
