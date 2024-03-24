package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"
)

// Create a Match struct
type Match struct {
	Status        string    `json:"status"`
	StartTime     time.Time `json:"start_time_iso"`
	HomeTeam      string    `json:"home_team"`
	HomeTeamScore int       `json:"home_team_score"`
	AwayTeam      string    `json:"away_team"`
	AwayTeamScore int       `json:"away_team_score"`
}

// Add a String method to the Match struct
func (m Match) String() string {
	return fmt.Sprintf(
		"%s - %-15s %2d - %-2d %-15s %-12s",
		m.StartTime.Local().Format(time.RFC1123),
		m.HomeTeam,
		m.HomeTeamScore,
		m.AwayTeamScore,
		m.AwayTeam,
		m.Status,
	)
}

// Create a Matches struct
type Matches struct {
	Matches []Match `json:"matches"`
}

// Add a String method to the Matches struct
func (m Matches) String() string {
	var str string

	// Loop through the matches and add them to the string
	for _, match := range m.Matches {
		str += match.String() + "\n"
	}

	// Trim the string
	str = strings.TrimSpace(str)

	// Return the string
	return str
}

type ctxKeys string

const (
	user ctxKeys = "user"
)

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

	// Redirect all requests without a trailing slash to one with a trailing slash
	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//     http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently)
	// })

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

func getMatches() (Matches, error) {
	// Matches struct
	var matches Matches

	// Get data from schleising.net/football/api/
	resp, err := http.Get("https://www.schleising.net/football/api/")

	// Check for errors
	if err != nil {
		return matches, err
	}

	// Check the status code
	if resp.StatusCode != 200 {
		return matches, err
	}

	// Defer the closing of the response body
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)

	// Check for errors
	if err != nil {
		return matches, err
	}

	// Unmarshal the response body into a Matches struct
	if err := json.Unmarshal(body, &matches); err != nil {
		return matches, err
	}

	return matches, nil
}
