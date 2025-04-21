package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux" // HTTP router with middleware support
	"github.com/joho/godotenv" // loads .env file into environment

	"github.com/soumyajitg16/eth-fundflow.git/pkg/analysis" // contains business logic for /beneficiary and /payer
	"github.com/soumyajitg16/eth-fundflow.git/pkg/client" // Etherscan API client
)

func main() {
	godotenv.Load() // reads .env into os.Environ()

	// Retrieve Etherscan API key from env
	apiKey := os.Getenv("ETHERSCAN_API_KEY")
	if apiKey == "" {
		log.Fatal("ETHERSCAN_API_KEY environment variable required") // internally calls os.exit()
	}

	// Initialize Etherscan client with timeout
	esClient := client.NewEtherscanClient(apiKey)

	// Set up router and attach middleware & handlers
	r := mux.NewRouter()
	r.Use(loggingMiddleware) // logs each incoming request
	r.HandleFunc("/beneficiary", analysis.BeneficiaryHandler(esClient)).Methods("GET")
	r.HandleFunc("/payer", analysis.PayerHandler(esClient)).Methods("GET")

	// Determine port (default 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// loggingMiddleware logs incoming HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}