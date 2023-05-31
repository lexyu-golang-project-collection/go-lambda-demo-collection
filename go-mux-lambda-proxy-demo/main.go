package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	gorillamuxAdapter "github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
)

var gorillaLambda *gorillamuxAdapter.GorillaMuxAdapter

/*
func init() {
	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Mux cold start")
	router := mux.NewRouter()
	router.HandleFunc("/ping", PingHandler).Methods("GET")
	router.HandleFunc("/hello", HelloHandler).Methods("GET")
	router.HandleFunc("/get_user", getUser).Methods("POST")

	gorillaLambda = gorillamux.New(router)
}
*/

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla pong!\n"))
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"name": "Lex", "email": "test@gmail.com"})
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla Hello, World\n"))
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (*core.SwitchableAPIGatewayResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return gorillaLambda.ProxyWithContext(ctx, *core.NewSwitchableAPIGatewayRequestV1(&req))
}

func main() {
	fmt.Printf("OS: %s\nArchitecture: %s\n", runtime.GOOS, runtime.GOARCH)

	log.Printf("Mux cold start")
	router := mux.NewRouter()
	router.HandleFunc("/ping", PingHandler).Methods("GET")
	router.HandleFunc("/hello", HelloHandler).Methods("GET")
	router.HandleFunc("/get_user", getUser).Methods("POST")

	env := os.Getenv("MUX_MODE")
	if env == "release" {
		gorillaLambda = gorillamuxAdapter.New(router)
		lambda.Start(Handler)
	} else {
		// Server
		server := http.Server{
			Addr:    ":8080",
			Handler: router,
		}
		// Run Server
		server.ListenAndServe()
	}
}
