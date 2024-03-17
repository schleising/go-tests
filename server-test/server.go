package main

import (
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

func main() {
	// Create a new server
	http.Handle("/api", http.HandlerFunc(handler))

	// Listen and serve on port 8080
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("%s %s %s\n", req.Method, req.URL.Path, req.Proto)
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
