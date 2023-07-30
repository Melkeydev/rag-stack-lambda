package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/golang-jwt/jwt"
)

type MyEvent struct {
	Name string `json:"name"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	userTableName = "user-table-name"
)

// Right now this doesnt do a whole much
// This was just to handle some logic

func RegisterHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var registerReq RegisterRequest

	// we need to unmarshal the request into our struct

	err := json.Unmarshal([]byte(request.Body), &registerReq)
	if err != nil {
		// I want to use a different loggign package
		log.Printf("Unable to unmarshal register request:%v", err)
		return events.APIGatewayProxyResponse{Body: "Invalid Request"}, err
	}

	fmt.Println(registerReq.Username)

	// We need to obviously validate the password and username
	// Create an AWS session and Insert these into my DynamoDB

	sess := session.Must(session.NewSession())
	DDB := dynamodb.New(sess)

	// Insert the item in dynamo
	item := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(registerReq.Username),
			},
			"password": {
				S: aws.String(registerReq.Password),
			},
		},
		TableName: aws.String(userTableName),
	}

	_, err = DDB.PutItem(item)
	if err != nil {
		log.Printf("Failed to input item into user DDB: %v", err)
		return events.APIGatewayProxyResponse{Body: "Internal Server Error - DDB", StatusCode: 500}, err
	}

	token, err := generateToken(registerReq.Username)
	if err != nil {
		log.Print("Could not issue jwt token")
		return events.APIGatewayProxyResponse{Body: "Internal Server Error - Generating token", StatusCode: 500}, err
	}

	responseBody := map[string]string{
		"token": token,
	}

	fmt.Println("this is response body", responseBody)

	responsejson, err := json.Marshal(responseBody)
	fmt.Println("this is response json", responsejson)

	if err != nil {
		log.Printf("Failed to marshal response %v", err)
	}

	return events.APIGatewayProxyResponse{Body: string(responsejson), StatusCode: 200}, nil

}

// Returns the actual token string and a error
func generateToken(username string) (string, error) {
	// TODO: this should come from env
	mySigningKey := []byte("randomString")

	type MyCustomClaims struct {
		Username string `json:"username"`
		jwt.StandardClaims
	}

	claims := MyCustomClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Printf("Failed to sign the token due to: %v", err)
		return "", err
	}

	return ss, nil

}

// I want login handler
// I want a register handler
// I also want a jwt function
func main() {
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/register":
			return RegisterHandler(request)
		default:
			return events.APIGatewayProxyResponse{Body: "Not Found", StatusCode: http.StatusNotFound}, nil
		}
	})
}
