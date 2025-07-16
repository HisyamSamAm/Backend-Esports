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

	fmt.Println("⚔️ Creating 12 teams (6 for each tournament, 5 players per team)...")

	ctx := context.Background()

	// Get all players to assign to teams
	cursor, err := config.PlayersCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to fetch players:", err)
	}
	defer cursor.Close(ctx)

	var players []model.Player
	if err = cursor.All(ctx, &players); err != nil {
		log.Fatal("Failed to decode players:", err)
	}

	if len(players) < 60 {
		log.Fatal("Need at least 60 players to create 12 teams. Please run seed-players first.")
	}

	fmt.Printf("📊 Found %d players in database\n", len(players))

	// Create teams with exactly 5 players each
	teams := []model.Team{
		// Tournament 1 Teams (players 0-29)
		{
			ID:        primitive.NewObjectID(),
			TeamName:  "RRQ Hoshi",
			CaptainID: players[0].ID,
			Members:   []primitive.ObjectID{players[0].ID, players[1].ID, players[2].ID, players[3].ID, players[4].ID},
			LogoURL:   "rrq_hoshi_logo.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			TeamName:  "EVOS Legends",
			CaptainID: players[5].ID,
			Members:   []primitive.ObjectID{players[5].ID, players[6].ID, players[7].ID, players[8].ID, players[9].ID},
			LogoURL:   "evos_legends_logo.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			TeamName:  "Bigetron Alpha",
			CaptainID: players[10].ID,
			Members:   []primitive.ObjectID{players[10].ID, players[11].ID, players[12].ID, players[13].ID, players[14].ID},
			LogoURL:   "bigetron_alpha_logo.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			TeamName:  "ONIC Esports",
			CaptainID: players[15].ID,
			Members:   []primitive.ObjectID{players[15].ID, players[16].ID, players[17].ID, players[18].ID, players[19].ID},
			LogoURL:   "onic_esports_logo.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			TeamName:  "Geek Fam",
			CaptainID: players[20].ID,
			Members:   []primitive.ObjectID{players[20].ID, players[21].ID, players[22].ID, players[23].ID, players[24].ID},
			LogoURL:   "geek_fam_logo.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			TeamName:  "Aura Fire",
			CaptainID: players[25].ID,
			Members:   []primitive.ObjectID{players[25].ID, players[26].ID, players[27].ID, players[28].ID, players[29].ID},
			LogoURL:   "aura_fire_logo.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		// Tournament 2 Teams (players 30-59)
		{
			ID:        primitive.NewObjectID(),
			TeamName:  "Rebellion Zion",
			CaptainID: players[30].ID,
			Members:   []primitive.ObjectID{players[30].ID, players[31].ID, players[32].ID, players[33].ID, players[34].ID},
			LogoURL:   "rebellion_zion_logo.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			TeamName:  "GPX Basreng",
			CaptainID: players[35].ID,
			Members:   []primitive.ObjectID{players[35].ID, players[36].ID, players[37].ID, players[38].ID, players[39].ID},
			LogoURL:   "gpx_basreng_logo.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			TeamName:  "Alter Ego",
			CaptainID: players[40].ID,
			Members:   []primitive.ObjectID{players[40].ID, players[41].ID, players[42].ID, players[43].ID, players[44].ID},
			LogoURL:   "alter_ego_logo.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			TeamName:  "Todak",
			CaptainID: players[45].ID,
			Members:   []primitive.ObjectID{players[45].ID, players[46].ID, players[47].ID, players[48].ID, players[49].ID},
			LogoURL:   "todak_logo.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			TeamName:  "Team Flash",
			CaptainID: players[50].ID,
			Members:   []primitive.ObjectID{players[50].ID, players[51].ID, players[52].ID, players[53].ID, players[54].ID},
			LogoURL:   "team_flash_logo.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        primitive.NewObjectID(),
			TeamName:  "Omega Esports",
			CaptainID: players[55].ID,
			Members:   []primitive.ObjectID{players[55].ID, players[56].ID, players[57].ID, players[58].ID, players[59].ID},
			LogoURL:   "omega_esports_logo.png",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Insert all teams
	var teamInterfaces []interface{}
	for _, team := range teams {
		teamInterfaces = append(teamInterfaces, team)
	}

	result, err := config.TeamsCollection.InsertMany(ctx, teamInterfaces)
	if err != nil {
		log.Fatal("Failed to insert teams:", err)
	}

	fmt.Printf("✅ Successfully created %d teams!\n", len(result.InsertedIDs))
	fmt.Println("\n📋 Teams breakdown:")
	fmt.Println("   🏆 Tournament 1 Teams (6 teams):")
	fmt.Println("     • RRQ Hoshi (Captain: Lemon)")
	fmt.Println("     • EVOS Legends (Captain: Wann)")
	fmt.Println("     • Bigetron Alpha (Captain: Branz)")
	fmt.Println("     • ONIC Esports (Captain: Sanz)")
	fmt.Println("     • Geek Fam (Captain: JessNoLimit)")
	fmt.Println("     • Aura Fire (Captain: Oura)")
	fmt.Println("   🏆 Tournament 2 Teams (6 teams):")
	fmt.Println("     • Rebellion Zion (Captain: Tuanmuda)")
	fmt.Println("     • GPX Basreng (Captain: Potato)")
	fmt.Println("     • Alter Ego (Captain: Ahmad)")
	fmt.Println("     • Todak (Captain: Moon)")
	fmt.Println("     • Team Flash (Captain: Flash)")
	fmt.Println("     • Omega Esports (Captain: Omega)")
	fmt.Println("\n✨ Each team has exactly 5 unique players with no overlaps!")
}
