package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	ragDynamo "ragStackECS/clients/dynamo"
	ragJWT "ragStackECS/clients/jwt"
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

func NewApp(db ragDynamo.UserStorageDB, jwt ragJWT.TokenValidator) *App {
	return &App{
		db:  ragDynamo.NewDynamoDBClient(),
		jwt: ragJWT.NewJWTClient(db),
	}
}

func (app *App) TestHandler(w http.ResponseWriter, r *http.Request) {

	responseBody := map[string]string{
		"message": "Hi you have hit this route",
	}
	responseJSON, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST, GET")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

}

func (app *App) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {

	responseBody := map[string]string{
		"message": "This is the health check for the server",
	}

	responseJSON, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST, GET")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

}

func main() {
	// Test out the RDS client
	db := ragDynamo.NewDynamoDBClient()
	jwt := ragJWT.NewJWTClient(db)

	app := NewApp(db, jwt)

	http.HandleFunc("/", app.HealthCheckHandler)
	http.HandleFunc("/test", app.TestHandler)

	port := ":8080"
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))

}
