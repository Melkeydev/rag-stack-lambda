package jwt

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt"
)

func ValidateJWTMiddleware(next func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		tokenString := extractTokenFromHeaders(request.Headers)
		if tokenString == "" {
			return events.APIGatewayProxyResponse{Body: "Missing Authentication Token", StatusCode: http.StatusUnauthorized}, nil
		}

		fmt.Println("this is the tokenString", tokenString)

		mySigningKey := []byte("randomString")

		token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})

		fmt.Println("this is the token", token)
		fmt.Println("this is the error", err)
		fmt.Println("this is if it is valid", token.Valid)

		if err != nil || !token.Valid {
			return events.APIGatewayProxyResponse{Body: "Invalid or Expired Token", StatusCode: http.StatusUnauthorized}, nil
		}

		claims, ok := token.Claims.(*MyCustomClaims)
		if !ok {
			return events.APIGatewayProxyResponse{Body: "Invalid Token Claims", StatusCode: http.StatusUnauthorized}, nil
		}

		request.RequestContext.Authorizer = map[string]interface{}{
			"username": claims.Username,
		}

		return next(request)
	}
}

func extractTokenFromHeaders(headers map[string]string) string {
	authHeader, ok := headers["Authorization"]
	if !ok {
		return ""
	}

	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return ""
	}

	return splitToken[1]
}
