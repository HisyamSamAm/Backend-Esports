package main

import (
	"context"
	"embeck/config"
	"embeck/model"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	fmt.Println("🏆 Creating 2 tournaments (6 teams each)...")

	ctx := context.Background()

	// Get all teams to assign to tournaments
	cursor, err := config.TeamsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to fetch teams:", err)
	}
	defer cursor.Close(ctx)

	var teams []model.Team
	if err = cursor.All(ctx, &teams); err != nil {
		log.Fatal("Failed to decode teams:", err)
	}

	if len(teams) < 12 {
		log.Fatal("Need at least 12 teams to create 2 tournaments. Please run seed-teams first.")
	}

	fmt.Printf("📊 Found %d teams in database\n", len(teams))

	// Create a dummy admin user ID for CreatedBy field
	adminID := primitive.NewObjectID()

	// Create 2 tournaments with 6 teams each
	tournaments := []model.Tournament{
		{
			ID:               primitive.NewObjectID(),
			Name:             "Mobile Legends Championship 2025",
			Description:      "The biggest Mobile Legends tournament in Indonesia featuring top professional teams competing for the ultimate prize.",
			StartDate:        time.Date(2025, 8, 15, 10, 0, 0, 0, time.UTC),
			EndDate:          time.Date(2025, 8, 25, 20, 0, 0, 0, time.UTC),
			PrizePool:        "Rp 2,500,000,000",
			RulesDocumentURL: "https://embeck.gg/tournaments/mlc2025/rules.pdf",
			Status:           "upcoming",
			TeamsParticipating: []primitive.ObjectID{
				teams[0].ID, // RRQ Hoshi
				teams[1].ID, // EVOS Legends
				teams[2].ID, // Bigetron Alpha
				teams[3].ID, // ONIC Esports
				teams[4].ID, // Geek Fam
				teams[5].ID, // Aura Fire
			},
			CreatedBy: adminID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:               primitive.NewObjectID(),
			Name:             "EMBECK Invitational Tournament 2025",
			Description:      "An exclusive invitational tournament featuring rising stars and veteran teams battling for glory and substantial prize money.",
			StartDate:        time.Date(2025, 9, 10, 9, 0, 0, 0, time.UTC),
			EndDate:          time.Date(2025, 9, 20, 18, 0, 0, 0, time.UTC),
			PrizePool:        "Rp 1,500,000,000",
			RulesDocumentURL: "https://embeck.gg/tournaments/embeck2025/rules.pdf",
			Status:           "upcoming",
			TeamsParticipating: []primitive.ObjectID{
				teams[6].ID,  // Rebellion Zion
				teams[7].ID,  // GPX Basreng
				teams[8].ID,  // Alter Ego
				teams[9].ID,  // Todak
				teams[10].ID, // Team Flash
				teams[11].ID, // Omega Esports
			},
			CreatedBy: adminID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Insert all tournaments
	var tournamentInterfaces []interface{}
	for _, tournament := range tournaments {
		tournamentInterfaces = append(tournamentInterfaces, tournament)
	}

	result, err := config.TournamentsCollection.InsertMany(ctx, tournamentInterfaces)
	if err != nil {
		log.Fatal("Failed to insert tournaments:", err)
	}

	fmt.Printf("✅ Successfully created %d tournaments!\n", len(result.InsertedIDs))
	fmt.Println("\n📋 Tournament details:")
	fmt.Println("   🏆 Mobile Legends Championship 2025:")
	fmt.Println("     • Date: August 15-25, 2025")
	fmt.Println("     • Prize Pool: Rp 2,500,000,000")
	fmt.Println("     • Teams: RRQ Hoshi, EVOS Legends, Bigetron Alpha, ONIC Esports, Geek Fam, Aura Fire")
	fmt.Println("   🏆 EMBECK Invitational Tournament 2025:")
	fmt.Println("     • Date: September 10-20, 2025")
	fmt.Println("     • Prize Pool: Rp 1,500,000,000")
	fmt.Println("     • Teams: Rebellion Zion, GPX Basreng, Alter Ego, Todak, Team Flash, Omega Esports")
	fmt.Println("\n✨ Each tournament has exactly 6 unique teams with no overlaps!")
	fmt.Println("📊 Summary: 2 tournaments × 6 teams × 5 players = 60 total unique players!")
}
