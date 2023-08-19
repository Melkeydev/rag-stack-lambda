package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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

type App struct {
	db  ragDynamo.UserStorageDB
	jwt ragJWT.TokenValidator
}

// TODO: is this even the best approach?
func NewApp(db ragDynamo.UserStorageDB, jwt ragJWT.TokenValidator) *App {
	return &App{
		db:  ragDynamo.NewDynamoDBClient(),
		jwt: ragJWT.NewJWTClient(db),
	}
}

func (app *App) RegisterHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var registerReq RegisterRequest

	err := json.Unmarshal([]byte(request.Body), &registerReq)
	if err != nil {
		log.Printf("Unable to unmarshal register request:%v", err)
		return events.APIGatewayProxyResponse{Body: "Invalid Request"}, err
	}

	// Validate if the user already exists in DynamoDB
	_, err = app.db.GetUser(registerReq.Username)
	if err == nil {
		return events.APIGatewayProxyResponse{Body: "Username already exists", StatusCode: http.StatusBadRequest}, nil
	}

	hashedPassword, err := ragCrypto.GeneratePassword(registerReq.Password)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Invalid request", StatusCode: http.StatusBadRequest}, nil
	}

	accessToken, err := app.jwt.GenerateAccessToken(registerReq.Username)
	if err != nil {
		log.Print("Could not issue jwt token")
		return events.APIGatewayProxyResponse{Body: "Internal Server Error - Generating token", StatusCode: 500}, err
	}

	refreshToken, err := app.jwt.GenerateRefreshToken(registerReq.Username)
	if err != nil {
		log.Print("Could not issue refresh token")
		return events.APIGatewayProxyResponse{Body: "Internal Server Error - Generating token", StatusCode: 500}, err
	}

	// Store the username and hashed password in DynamoDB
	err = app.db.AddUserToDB(registerReq.Username, string(hashedPassword), refreshToken)
	if err != nil {
		log.Printf("Failed to add user to DynamoDB: %v", err)
		return events.APIGatewayProxyResponse{Body: "Internal Server Error - DDB", StatusCode: http.StatusInternalServerError}, nil
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		Path:     "/",
	}

	response := events.APIGatewayProxyResponse{
		Body:              fmt.Sprintf(`{"access_token": "%s"}`, accessToken),
		StatusCode:        http.StatusOK,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string{"Set-Cookie": {refreshCookie.String()}},
	}

	return response, nil
}

func (app *App) LoginHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var loginReq LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginReq)
	if err != nil {
		log.Printf("Unable to unmarshal login request:%v", err)
		return events.APIGatewayProxyResponse{Body: "Invalid Request"}, err
	}

	// Check if the user exists in DynamoDB
	user, err := app.db.GetUser(loginReq.Username)
	if err != nil {
		log.Printf("Failed to retrieve user from DynamoDB: %v", err)
		return events.APIGatewayProxyResponse{Body: "Invalid username or password", StatusCode: http.StatusUnauthorized}, nil
	}

	// Check password is correct for user
	if !ragCrypto.ComparePasswords(user.Password, loginReq.Password) {
		return events.APIGatewayProxyResponse{Body: "Invalid username or password", StatusCode: http.StatusUnauthorized}, nil
	}

	// Create new access token
	accessToken, err := app.jwt.GenerateAccessToken(loginReq.Username)
	if err != nil {
		log.Print("Could not issue jwt token")
		return events.APIGatewayProxyResponse{Body: "Internal Server Error - Generating token", StatusCode: 500}, err
	}

	// create new refresh token
	refreshToken, err := app.jwt.GenerateRefreshToken(loginReq.Username)
	if err != nil {
		log.Print("Could not issue refresh token")
		return events.APIGatewayProxyResponse{Body: "Internal Server Error - Generating token", StatusCode: 500}, err
	}

	// update refresh token in DynamoDB
	err = app.db.UpdateUserToken(loginReq.Username, refreshToken)
	if err != nil {
		log.Printf("Failed to update user's token in DynamoDB: %v", err)
		return events.APIGatewayProxyResponse{Body: "Internal Server Error", StatusCode: http.StatusInternalServerError}, nil
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		Path:     "/",
	}

	response := events.APIGatewayProxyResponse{
		Body:              fmt.Sprintf(`{"access_token": "%s"}`, accessToken),
		StatusCode:        http.StatusOK,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string{"Set-Cookie": {refreshCookie.String()}},
	}

	return response, nil
}

func (app *App) ProtectedHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	username, ok := request.RequestContext.Authorizer["username"].(string)
	if !ok {
		return events.APIGatewayProxyResponse{Body: "Unauthorized", StatusCode: http.StatusForbidden}, nil
	}

	responseBody := fmt.Sprintf("Hey %s - this is a protected route", username)

	return events.APIGatewayProxyResponse{Body: string(responseBody), StatusCode: 200}, nil
}

func extractRefreshTokenFromCookie(request events.APIGatewayProxyRequest) (string, error) {
	cookieHeader := request.Headers["Cookie"]
	cookies := strings.Split(cookieHeader, "; ")
	for _, cookie := range cookies {
		parts := strings.SplitN(cookie, "=", 2)
		if len(parts) == 2 && parts[0] == "refresh_token" {
			return parts[1], nil
		}
	}
	return "", errors.New("Refresh Token not found in cookies")
}

func (app *App) RefreshHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	refreshToken, err := extractRefreshTokenFromCookie(request)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Missing Refresh Token", StatusCode: http.StatusBadRequest}, nil
	}

	username, err := app.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Invalid Refresh Token", StatusCode: http.StatusUnauthorized}, nil
	}

	accessToken, err := app.jwt.GenerateAccessToken(username)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Internal Server Error", StatusCode: http.StatusInternalServerError}, nil
	}

	// Create an HTTP-only cookie for the new refresh token
	newRefreshToken, err := app.jwt.GenerateRefreshToken(username)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Internal Server Error", StatusCode: http.StatusInternalServerError}, nil
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		Path:     "/",
	}

	// Create a response with the new access token and set the refresh cookie
	response := events.APIGatewayProxyResponse{
		Body:              fmt.Sprintf(`{"access_token": "%s"}`, accessToken),
		StatusCode:        http.StatusOK,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string{"Set-Cookie": {refreshCookie.String()}},
	}

	return response, nil
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
		case "/protected":
			return ragJWT.ValidateJWTMiddleware(app.ProtectedHandler)(request)
		default:
			return events.APIGatewayProxyResponse{Body: "Not Found", StatusCode: http.StatusNotFound}, nil
		}
	})
}
