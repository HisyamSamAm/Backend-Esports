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

	fmt.Println("🧹 Cleaning up all existing data...")

	ctx := context.Background()

	// Delete in order to avoid foreign key issues
	fmt.Println("1. Deleting tickets...")
	ticketsResult, err := config.TicketsCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to delete tickets:", err)
	}
	fmt.Printf("✅ Deleted %d tickets\n", ticketsResult.DeletedCount)

	fmt.Println("2. Deleting matches...")
	matchesResult, err := config.MatchesCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to delete matches:", err)
	}
	fmt.Printf("✅ Deleted %d matches\n", matchesResult.DeletedCount)

	fmt.Println("3. Deleting tournaments...")
	tournamentsResult, err := config.TournamentsCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to delete tournaments:", err)
	}
	fmt.Printf("✅ Deleted %d tournaments\n", tournamentsResult.DeletedCount)

	fmt.Println("4. Deleting teams...")
	teamsResult, err := config.TeamsCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to delete teams:", err)
	}
	fmt.Printf("✅ Deleted %d teams\n", teamsResult.DeletedCount)

	fmt.Println("5. Deleting players...")
	playersResult, err := config.PlayersCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to delete players:", err)
	}
	fmt.Printf("✅ Deleted %d players\n", playersResult.DeletedCount)

	fmt.Println("6. Deleting users...")
	usersResult, err := config.UsersCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to delete users:", err)
	}
	fmt.Printf("✅ Deleted %d users\n", usersResult.DeletedCount)

	fmt.Println("\n🎉 Database cleanup completed!")
}
