package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Match struct {
	Status string `json:"status"`
	StartTime time.Time `json:"start_time_iso"`
	HomeTeam string `json:"home_team"`
	HomeTeamScore int `json:"home_team_score"`
	AwayTeam string `json:"away_team"`
	AwayTeamScore int `json:"away_team_score"`
}

type Matches struct {
	Matches []Match `json:"matches"`
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
	err = json.Unmarshal(body, &matches)

	// Check for errors
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the matches
	for _, match := range matches.Matches {
		fmt.Printf(
			"%s - %-15s %2d - %-2d %-15s %-12s\n",
			match.StartTime.Local().Format(time.RFC1123),
			match.HomeTeam,
			match.HomeTeamScore,
			match.AwayTeamScore,
			match.AwayTeam,
			match.Status,
		)
	}
}
