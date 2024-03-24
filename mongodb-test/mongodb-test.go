package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"mongodb-test/models"
)

var logger *log.Logger

func main() {
	// Set up logging
	logger.Println("Starting MongoDB test")

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://macmini2:27017"))

	// Check for errors
	if err != nil {
		logger.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Ping the MongoDB server
	err = client.Ping(ctx, readpref.Primary())

	// Check for errors
	if err != nil {
		logger.Fatalf("Error pinging MongoDB: %v", err)
	}

	// Log success
	logger.Println("Connected to MongoDB")

	// Get the web_database
	database := client.Database("web_database")

	// Get the pl_matches_2023_2024 collection
	collection := database.Collection("pl_matches_2023_2024")

	// Get the Liverpool vs Manchester City match
	match := getOneMatch(ctx, collection, "Liverpool", "Man City")

	// Log the match
	logger.Printf("%v", match)

	// Get all Liverpool matches
	matches := getAllTeamMatches(ctx, collection, "Liverpool")

	// Log the matches
	logger.Printf("%v", matches)
}

func getOneMatch(ctx context.Context, collection *mongo.Collection, homeTeam, awayTeam string) models.Match {
	// Get the document that has homeTeam as the home team and awayTeam as the away team
	var result models.Match
	err := collection.FindOne(
		ctx,
		bson.D{
			{Key: "home_team.short_name", Value: homeTeam},
			{Key: "away_team.short_name", Value: awayTeam},
		},
	).Decode(&result)

	// Check for errors
	if err != nil {
		logger.Fatalf("Error getting first document: %v", err)
	}

	// Return the result
	return result
}

func getAllTeamMatches(ctx context.Context, collection *mongo.Collection, team string) models.MatchList {
	// Get all documents that have `team` as the home team or away team
	cursor, err := collection.Find(
		ctx,
		bson.D{
			{Key: "$or", Value: bson.A{
				bson.D{{Key: "home_team.short_name", Value: team}},
				bson.D{{Key: "away_team.short_name", Value: team}},
			}},
		},
	)

	// Check for errors
	if err != nil {
		logger.Fatalf("Error getting documents: %v", err)
	}

	// Iterate through the cursor
	defer cursor.Close(ctx)

	// Create a MatchList
	var matchList models.MatchList

	// Iterate through the cursor
	for cursor.Next(ctx) {
		// Decode the document
		var result models.Match
		err := cursor.Decode(&result)

		// Check for errors
		if err != nil {
			logger.Fatalf("Error decoding document: %v", err)
		}

		// Append the result to the matchList
		matchList.Matches = append(matchList.Matches, result)
	}

	// Return the matchList
	return matchList
}

func init() {
	// Set up logging
	logger = log.New(log.Writer(), "mongodb-test: ", log.Ldate|log.Ltime|log.Lshortfile)
}
