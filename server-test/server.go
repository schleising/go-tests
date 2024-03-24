package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	// Error variable
	var err error

	// Create a new MongoDB test
	mongo, err = mongodb_test.NewMongoTest()

	// Check for errors
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}

	// Create a new channel
	c := make(chan os.Signal, 1)

	// Notify the channel of the following signals
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// Create a goroutine to listen for signals
	go func() {
		// Wait for a signal
		sig := <-c

		// Print the signal
		fmt.Println("Received signal:", sig)

		// Close the MongoDB connection
		mongo.Close()

		// Exit the program
		os.Exit(0)
	}()

	// Listen and serve on port 8080
	http.ListenAndServe(":8080", handler)
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
