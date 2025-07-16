package config

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database configuration
var DBName = "turnamen_esport"
var MongoString string = os.Getenv("MONGOSTRING")

// Global database instance
var DB *mongo.Database

// Collection instances
var UsersCollection *mongo.Collection
var TournamentsCollection *mongo.Collection
var TeamsCollection *mongo.Collection
var PlayersCollection *mongo.Collection
var MatchesCollection *mongo.Collection
var TicketsCollection *mongo.Collection
var TransactionsCollection *mongo.Collection

// MongoConnect establishes connection to MongoDB and returns database instance
func MongoConnect(dbname string) (db *mongo.Database) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MongoString))
	if err != nil {
		fmt.Printf("MongoConnect: %v\n", err)
		return nil
	}

	// Test the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Printf("MongoDB Ping failed: %v\n", err)
		return nil
	}

	fmt.Printf("✅ Successfully connected to MongoDB database: %s\n", dbname)

	// Set global database instance
	DB = client.Database(dbname)

	// Initialize collections
	UsersCollection = DB.Collection("users")
	TournamentsCollection = DB.Collection("tournaments")
	TeamsCollection = DB.Collection("teams")
	PlayersCollection = DB.Collection("players")
	MatchesCollection = DB.Collection("matches")
	TicketsCollection = DB.Collection("tickets")
	TransactionsCollection = DB.Collection("transactions")

	return DB
}
