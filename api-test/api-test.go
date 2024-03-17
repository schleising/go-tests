package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	// Get data from schleising.net/football/api/
	resp, err := http.Get("https://www.schleising.net/football/api/")

	// Check for errors
	if err != nil {
		fmt.Println(err)
		return
	}

	// Check the status code
	if resp.StatusCode != 200 {
		fmt.Println("Error: ", resp.Status)
		return
	}

	// Defer the closing of the response body
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)

	// Check for errors
	if err != nil {
		fmt.Println(err)
		return
	}

	// Unmarshal the response body into a Matches struct
	var matches Matches
	if err := json.Unmarshal(body, &matches); err != nil {
		fmt.Println(err)
		return
	}

	// Print the matches from the API
	fmt.Println()
	fmt.Println(("Matches from the API"))
	fmt.Println(matches)

	// Marshall the matches back to JSON and save it to a file
	matchesJSON, err := json.MarshalIndent(matches, "", "  ")

	// Check for errors
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set the file path
	file := "tests/matches.json"

	// Create a tests directory
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		fmt.Println(err)
		return
	}

	// Write the JSON to a file
	if err := os.WriteFile(file, matchesJSON, 0644); err != nil {
		fmt.Println(err)
		return
	}

	// Read the JSON from the file
	bytes, err := os.ReadFile(file)

	// Check for errors
	if err != nil {
		fmt.Println(err)
		return
	}

	// Unmarshal the JSON into a Matches struct
	var matchesFromFile Matches
	if err := json.Unmarshal(bytes, &matchesFromFile); err != nil {
		fmt.Println(err)
		return
	}

	// Print the matches from the file
	fmt.Println()
	fmt.Println(("Matches from the file"))
	fmt.Println(matchesFromFile)
}
