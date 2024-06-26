package mongodb_test

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

type MongoTest struct {
	client *mongo.Client
	database *mongo.Database
	collection *mongo.Collection
	logger *log.Logger
}

func NewMongoTest() (*MongoTest, error) {
	// Set up logging
	logger := log.New(log.Writer(), "mongodb-test: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Set up logging
	logger.Println("Starting MongoDB test")

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://macmini2:27017"))

	// Check for errors
	if err != nil {
		// Log the error
		logger.Printf("Error connecting to MongoDB: %v", err)

		// Return an error
		return nil, err
	}

	// Ping the MongoDB server
	err = client.Ping(ctx, readpref.Primary())

	// Check for errors
	if err != nil {
		// Log the error
		logger.Printf("Error pinging MongoDB: %v", err)

		// Return an error
		return nil, err
	}

	// Log success
	logger.Println("Connected to MongoDB")

	// Get the web_database
	database := client.Database("web_database")

	// Get the pl_matches_2023_2024 collection
	collection := database.Collection("pl_matches_2023_2024")

	// Create a new MongoTest struct
	m := &MongoTest{
		client: client,
		database: database,
		collection: collection,
		logger: logger,
	}

	// Return the MongoTest struct
	return m, nil
}

func (m *MongoTest) Close() {
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Disconnect from MongoDB
	err := m.client.Disconnect(ctx)

	// Check for errors
	if err != nil {
		m.logger.Fatalf("Error disconnecting from MongoDB: %v", err)
	}

	// Log success
	m.logger.Println("Disconnected from MongoDB")
}

func (m *MongoTest) GetOneMatch(homeTeam, awayTeam string) (models.Match, error) {
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Get the document that has homeTeam as the home team and awayTeam as the away team
	var result models.Match
	err := m.collection.FindOne(
		ctx,
		bson.D{
			{Key: "home_team.short_name", Value: homeTeam},
			{Key: "away_team.short_name", Value: awayTeam},
		},
	).Decode(&result)

	// Check for errors
	if err != nil {
		m.logger.Printf("Error getting first document: %v", err)
		return models.Match{}, err
	}

	// Return the result
	return result, nil
}

func (m *MongoTest) GetAllTeamMatches(team string) (models.MatchList, error) {
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Get all documents that have `team` as the home team or away team
	cursor, err := m.collection.Find(
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
		// Log the error
		m.logger.Printf("Error getting all documents: %v", err)

		// Return an empty MatchList and the error
		return models.MatchList{}, err
	}

	// Close the cursor when the function returns
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
			// Log the error
			m.logger.Printf("Error decoding document: %v", err)

			// Return an empty MatchList and the error
			return models.MatchList{}, err
		}

		// Append the result to the matchList
		matchList.Matches = append(matchList.Matches, result)
	}

	// Return the matchList
	return matchList, nil
}

func (m *MongoTest) GetTodaysMatches() (models.MatchList, error) {
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Get the start of today
	start_date := time.Now().Truncate(24 * time.Hour)

	// Get the end of today
	end_date := start_date.Add(24 * time.Hour)
	
	// Get all documents between start_date and end_date
	cursor, err := m.collection.Find(
		ctx,
		bson.D{
			{Key: "utc_date", Value: bson.D{
				{Key: "$gte", Value: start_date},
				{Key: "$lt", Value: end_date},
			}},
		},
	)

	// Check for errors
	if err != nil {
		// Log the error
		m.logger.Printf("Error getting all documents: %v", err)

		// Return an empty MatchList and the error
		return models.MatchList{}, err
	}

	// Close the cursor when the function returns
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
			// Log the error
			m.logger.Printf("Error decoding document: %v", err)

			// Return an empty MatchList and the error
			return models.MatchList{}, err
		}

		// Append the result to the matchList
		matchList.Matches = append(matchList.Matches, result)
	}

	// Return the matchList
	return matchList, nil
}
