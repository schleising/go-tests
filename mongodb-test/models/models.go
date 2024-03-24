package models

import (
	"fmt"
	"time"
)

type Area struct {
	Id   int    `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
	Code string `bson:"code" json:"code"`
	Flag string `bson:"flag" json:"flag"`
}

type Competition struct {
	Id     int    `bson:"id" json:"id"`
	Name   string `bson:"name" json:"name"`
	Code   string `bson:"code" json:"code"`
	Type   string `bson:"type" json:"type"`
	Emblem string `bson:"emblem" json:"emblem"`
}

type Team struct {
	Id        int    `bson:"id" json:"id"`
	Name      string `bson:"name" json:"name"`
	ShortName string `bson:"short_name" json:"shortName"`
	Tla       string `bson:"tla" json:"tla"`
	Crest     string `bson:"crest" json:"crest"`
}

func (t Team) String() string {
	switch t.ShortName {
	case "Brighton Hove":
		return "Brighton"
	case "Wolverhampton":
		return "Wolves"
	case "Nottingham":
		return "Notts Forest"
	default:
		return t.ShortName
	}
}

type Season struct {
	Id              int    `bson:"id" json:"id"`
	StartDate       string `bson:"start_date" json:"startDate"`
	EndDate         string `bson:"end_date" json:"endDate"`
	CurrentMatchday int    `bson:"current_matchday" json:"currentMatchday"`
	Winner          Team   `bson:"winner" json:"winner"`
}

type MatchStatus string

const (
	Scheduled MatchStatus = "SCHEDULED"
	Timed     MatchStatus = "TIMED"
	InPlay    MatchStatus = "IN_PLAY"
	Paused    MatchStatus = "PAUSED"
	Finished  MatchStatus = "FINISHED"
	Suspended MatchStatus = "SUSPENDED"
	Postponed MatchStatus = "POSTPONED"
	Cancelled MatchStatus = "CANCELLED"
	Awarded   MatchStatus = "AWARDED"
)

func (m MatchStatus) String() string {
	switch m {
	case Scheduled, Timed, Awarded:
		return "Scheduled"
	case InPlay:
		return "In Play"
	case Paused:
		return "Half Time"
	case Finished:
		return "Full Time"
	case Suspended:
		return "Suspended"
	case Postponed:
		return "Postponed"
	case Cancelled:
		return "Cancelled"
	default:
		return "Unknown"
	}
}

type FullTime struct {
	Home int `bson:"home" json:"home"`
	Away int `bson:"away" json:"away"`
}

type HalfTime struct {
	Home int `bson:"home" json:"home"`
	Away int `bson:"away" json:"away"`
}

type Score struct {
	Winner   string   `bson:"winner" json:"winner"`
	Duration string   `bson:"duration" json:"duration"`
	FullTime FullTime `bson:"full_time" json:"fullTime"`
	HalfTime HalfTime `bson:"half_time" json:"halfTime"`
}

type Odds struct {
	Msg string `bson:"msg" json:"msg"`
}

type Referee struct {
	Id          int    `bson:"id" json:"id"`
	Name        string `bson:"name" json:"name"`
	Type        string `bson:"type" json:"type"`
	Nationality string `bson:"nationality" json:"nationality"`
}

type Match struct {
	Area        Area        `bson:"area" json:"area"`
	Competition Competition `bson:"competition" json:"competition"`
	Season      Season      `bson:"season" json:"season"`
	Id          int         `bson:"id" json:"id"`
	UtcDate     time.Time   `bson:"utc_date" json:"utcDate"`
	Status      MatchStatus `bson:"status" json:"status"`
	Matchday    int         `bson:"matchday" json:"matchday"`
	Stage       string      `bson:"stage" json:"stage"`
	Group       string      `bson:"group" json:"group"`
	LastUpdated time.Time   `bson:"last_updated" json:"lastUpdated"`
	HomeTeam    Team        `bson:"home_team" json:"homeTeam"`
	AwayTeam    Team        `bson:"away_team" json:"awayTeam"`
	Score       Score       `bson:"score" json:"score"`
	Odds        Odds        `bson:"odds" json:"odds"`
	Referees    []Referee   `bson:"referees" json:"referees"`
}

func (m Match) String() string {
	return fmt.Sprintf(
		"%v %14v %v - %v %-14v %v",
		m.UtcDate.Local().Format(time.RFC822),
		m.HomeTeam,
		m.Score.FullTime.Home,
		m.Score.FullTime.Away,
		m.AwayTeam,
		m.Status,
	)
}

type MatchList struct {
	Matches []Match `bson:"matches" json:"matches"`
}

func (ml MatchList) String() string {
	var s string
	s += fmt.Sprintln()
	for _, m := range ml.Matches {
		s += fmt.Sprintf("%v\n", m)
	}
	return s
}
