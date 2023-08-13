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
		// the actual logic for our middleware
		tokenString := extractTokenFromHeader(request.Headers)
		if tokenString == "" {
			return events.APIGatewayProxyResponse{Body: "Missing Auth Token", StatusCode: http.StatusUnauthorized}, nil
		}

		fmt.Println("this is the tokenString", tokenString)

		mySigningKey := []byte("randomString")

		token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})

		fmt.Println("this is the token", token)
		fmt.Println("this is the error", err)

		if err != nil || !token.Valid {
			return events.APIGatewayProxyResponse{Body: "Invalid or Expired Token", StatusCode: http.StatusForbidden}, nil
		}

		claims, ok := token.Claims.(*MyCustomClaims)
		if !ok {
			return events.APIGatewayProxyResponse{Body: "Invalid Token or Expired Token", StatusCode: http.StatusForbidden}, nil
		}

		request.RequestContext.Authorizer = map[string]interface{}{
			"username": claims.Username,
		}

		// once we pass all the logic we will kick off the next function in the chain
		return next(request)
	}
}

// We need to use this inside our middleware when we get the request
func extractTokenFromHeader(headers map[string]string) string {
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
