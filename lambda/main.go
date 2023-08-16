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

func RegisterHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var registerReq RegisterRequest

	err := json.Unmarshal([]byte(request.Body), &registerReq)
	if err != nil {
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

	accessToken, err := ragJWT.GenerateAccessToken(registerReq.Username)
	if err != nil {
		log.Print("Could not issue jwt token")
		return events.APIGatewayProxyResponse{Body: "Internal Server Error - Generating token", StatusCode: 500}, err
	}

	refreshToken, err := ragJWT.GenerateRefreshToken(registerReq.Username)
	if err != nil {
		log.Print("Could not issue refresh token")
		return events.APIGatewayProxyResponse{Body: "Internal Server Error - Generating token", StatusCode: 500}, err
	}

	// Store the username and hashed password in DynamoDB
	err = ragDynamo.AddUserToDynamoDB(registerReq.Username, string(hashedPassword), refreshToken)
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

func LoginHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var loginReq LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginReq)
	if err != nil {
		log.Printf("Unable to unmarshal login request:%v", err)
		return events.APIGatewayProxyResponse{Body: "Invalid Request"}, err
	}

	user, err := ragDynamo.GetUserFromDynamoDB(loginReq.Username)
	if err != nil {
		log.Printf("Failed to retrieve user from DynamoDB: %v", err)
		return events.APIGatewayProxyResponse{Body: "Invalid username or password", StatusCode: http.StatusUnauthorized}, nil
	}

	if !ragCrypto.ComparePasswords(user.Password, loginReq.Password) {
		return events.APIGatewayProxyResponse{Body: "Invalid username or password", StatusCode: http.StatusUnauthorized}, nil
	}

	accessToken, err := ragJWT.GenerateAccessToken(loginReq.Username)
	if err != nil {
		log.Print("Could not issue jwt token")
		return events.APIGatewayProxyResponse{Body: "Internal Server Error - Generating token", StatusCode: 500}, err
	}

	refreshToken, err := ragJWT.GenerateRefreshToken(loginReq.Username)
	if err != nil {
		log.Print("Could not issue refresh token")
		return events.APIGatewayProxyResponse{Body: "Internal Server Error - Generating token", StatusCode: 500}, err
	}

	err = ragDynamo.UpdateUserTokenInDynamoDB(loginReq.Username, refreshToken)
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

func ProtectedHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

func RefreshHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	refreshToken, err := extractRefreshTokenFromCookie(request)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Missing Refresh Token", StatusCode: http.StatusBadRequest}, nil
	}

	username, err := ragJWT.ValidateRefreshToken(refreshToken)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Invalid Refresh Token", StatusCode: http.StatusUnauthorized}, nil
	}

	accessToken, err := ragJWT.GenerateAccessToken(username)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Internal Server Error", StatusCode: http.StatusInternalServerError}, nil
	}

	// Create an HTTP-only cookie for the new refresh token
	newRefreshToken, err := ragJWT.GenerateRefreshToken(username)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Internal Server Error", StatusCode: http.StatusInternalServerError}, nil
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HttpOnly: true,
		Secure:   true,                                // Ensure your site uses HTTPS for this to work
		SameSite: http.SameSiteLaxMode,                // Lax, Strict, or None
		Expires:  time.Now().Add(30 * 24 * time.Hour), // 30 days from now
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
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/login":
			return LoginHandler(request)
		case "/register":
			return RegisterHandler(request)
		case "/refresh":
			return RefreshHandler(request)
		case "/protected":
			return ragJWT.ValidateJWTMiddleware(ProtectedHandler)(request)
		default:
			return events.APIGatewayProxyResponse{Body: "Not Found", StatusCode: http.StatusNotFound}, nil
		}
	})
}
