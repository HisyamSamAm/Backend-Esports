package main

import (
	"context"
	"embeck/config"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	if _, err := os.Stat("../../.env"); err == nil {
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Println("Warning: Could not load .env file")
		}
	}
}

func main() {
	db := config.MongoConnect(config.DBName)
	if db == nil {
		log.Fatal("Failed to connect to database")
	}

	fmt.Println("🔍 Verifying seeded data...")
	ctx := context.Background()

	playerCount, err := config.PlayersCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to count players:", err)
	}

	teamCount, err := config.TeamsCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to count teams:", err)
	}

	tournamentCount, err := config.TournamentsCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to count tournaments:", err)
	}

	matchCount, err := config.MatchesCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to count matches:", err)
	}

	ticketCount, err := config.TicketsCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to count tickets:", err)
	}

	userCount, err := config.UsersCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to count users:", err)
	}

	fmt.Printf("📊 Database Statistics:\n")
	fmt.Printf("   Players: %d\n", playerCount)
	fmt.Printf("   Teams: %d\n", teamCount)
	fmt.Printf("   Tournaments: %d\n", tournamentCount)
	fmt.Printf("   Matches: %d\n", matchCount)
	fmt.Printf("   Tickets: %d\n", ticketCount)
	fmt.Printf("   Users: %d\n", userCount)

	fmt.Println("\n🎉 Data verification completed!")
}
