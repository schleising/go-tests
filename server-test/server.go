package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mongodb-test"
	"mongodb-test/models"
)

type ctxKeys string

const (
	user ctxKeys = "user"
)

var mongo *mongodb_test.MongoTest

func main() {
	// Create a new mux
	mux := http.NewServeMux()

	// Handle the root route
	mux.Handle("/", http.HandlerFunc(rootHandler))

	// Handle the /api/ route
	mux.Handle("/api/", http.HandlerFunc(apiHandler))

	// Handle the /api/go/ route
	mux.Handle("/api/go/", http.HandlerFunc(goHandler))

	// Wrap the handlers with the logHandler function
	handler := logHandler(mux)

	// Create a new server
	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	// Creare an error variable
	var err error

	// Create a new MongoDB test
	mongo, err = mongodb_test.NewMongoTest()

	// Check for errors
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}

	// Close the MongoDB connection
	defer mongo.Close()

	// Create a goroutine to listen for signals
	go func() {
		// Start the server
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Println("Server Error:", err)
		}

		// Print that the server has stopped
		fmt.Println("Server stopped listening")
	}()

	// Create a new channel
	c := make(chan os.Signal, 1)

	// Notify the channel of the following signals
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a signal
	sig := <-c

	// Print the signal
	fmt.Println("Received Signal:", sig)

	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Error shutting down server:", err)
	}

	// Print that the server has stopped
	fmt.Println("Server stopped gracefully")
}

// Function to wrap the handlers and print the function name, request method, URL path, and protocol
func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Print the calling function name, request method, URL path, and protocol
		fmt.Printf("%s %s %s\n", req.Method, req.URL.Path, req.Proto)

		// Add some data to the request context
		ctx := req.Context()
		ctx = context.WithValue(ctx, user, "steve")
		req = req.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(w, req)

		// Write the function name, request method, URL path, and protocol to the response
		fmt.Fprintf(w, "Request: %s %s %s\n", req.Method, req.URL.Path, req.Proto)

		// Print the query string
		query := req.URL.Query()

		// Print the query string
		if len(query) == 0 {
			fmt.Println("No query string")
		} else {
			fmt.Println("Query string:")
		}

		// Loop through the query string
		for key, value := range query {
			fmt.Printf("%s\n", key)
			for _, v := range value {
				fmt.Printf("  %s\n", v)
			}
		}
	})
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	// Write the response
	fmt.Println("Root handler")

	// Print the data from the request context
	ctx := req.Context()
	value := ctx.Value(user)
	fmt.Println("User:", value)
}

func goHandler(w http.ResponseWriter, req *http.Request) {
	// Write the response
	fmt.Println("Go handler")

	// Print the data from the request context
	ctx := req.Context()
	value := ctx.Value(user)
	fmt.Println("User:", value)
}

func apiHandler(w http.ResponseWriter, req *http.Request) {
	// Check that the request is for the /api/ route
	if req.URL.Path != "/api/" {
		http.NotFound(w, req)
		return
	}

	fmt.Println("API handler")
	// Get the matches
	matches, err := getMatches()

	// Check for errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new template
	tmpl, err := template.ParseFiles("template.html")

	// Check for errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse the template
	if err := tmpl.Execute(w, matches); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getMatches() (models.MatchList, error) {
	// Get the matches
	matches, err := mongo.GetAllTeamMatches("Liverpool")

	// Check for errors
	if err != nil {
		return models.MatchList{}, err
	}

	// Return the matches
	return matches, nil
}
