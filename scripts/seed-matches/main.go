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
	// Load .env file if exists
	if _, err := os.Stat("../../.env"); err == nil {
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Println("Warning: Could not load .env file")
		}
	}
}

func main() {
	// Connect to database
	db := config.MongoConnect(config.DBName)
	if db == nil {
		log.Fatal("Failed to connect to database")
	}

	fmt.Println("🏆 Creating sample matches for tournaments...")

	ctx := context.Background()

	// Get existing tournaments to create matches for
	var tournaments []model.Tournament
	cursor, err := config.TournamentsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to fetch tournaments:", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &tournaments); err != nil {
		log.Fatal("Failed to decode tournaments:", err)
	}

	if len(tournaments) == 0 {
		log.Fatal("No tournaments found. Please run seed-tournaments first.")
	}

	// Get existing teams to create matches with
	var teams []model.Team
	teamsCursor, err := config.TeamsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to fetch teams:", err)
	}
	defer teamsCursor.Close(ctx)

	if err := teamsCursor.All(ctx, &teams); err != nil {
		log.Fatal("Failed to decode teams:", err)
	}

	if len(teams) < 2 {
		log.Fatal("Need at least 2 teams to create matches. Please run seed-teams first.")
	}

	// Clear existing matches
	_, err = config.MatchesCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to clear existing matches:", err)
	}

	// Create sample matches for each tournament
	var allMatches []interface{}
	matchCount := 0

	for _, tournament := range tournaments {
		// Create matches between participating teams
		participatingTeams := tournament.TeamsParticipating
		if len(participatingTeams) < 2 {
			// If no participating teams, use first few teams
			for i := 0; i < min(4, len(teams)); i++ {
				participatingTeams = append(participatingTeams, teams[i].ID)
			}
		}

		// Create group stage matches
		for i := 0; i < len(participatingTeams); i++ {
			for j := i + 1; j < len(participatingTeams); j++ {
				match := model.Match{
					ID:           primitive.NewObjectID(),
					TournamentID: tournament.ID,
					TeamAID:      participatingTeams[i],
					TeamBID:      participatingTeams[j],
					MatchDate:    tournament.StartDate.Add(time.Duration(matchCount*2) * time.Hour * 24),
					MatchTime:    "19:00 WIB",
					Location:     "JCC Senayan",
					Round:        "Group Stage",
					Status:       "scheduled",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}

				// Add some completed matches with scores
				if matchCount%3 == 0 {
					match.Status = "completed"
					scoreA := 2
					scoreB := 1
					match.ResultTeamAScore = &scoreA
					match.ResultTeamBScore = &scoreB
					match.WinnerTeamID = &participatingTeams[i]
					match.MatchDate = tournament.StartDate.Add(-time.Duration(matchCount) * time.Hour * 24)
				} else if matchCount%5 == 0 {
					match.Status = "ongoing"
					match.MatchDate = time.Now()
				}

				allMatches = append(allMatches, match)
				matchCount++

				// Limit matches per tournament
				if matchCount >= 10 {
					break
				}
			}
			if matchCount >= 10 {
				break
			}
		}

		// Add playoff matches if we have enough teams
		if len(participatingTeams) >= 4 {
			// Semi-final 1
			semifinal1 := model.Match{
				ID:           primitive.NewObjectID(),
				TournamentID: tournament.ID,
				TeamAID:      participatingTeams[0],
				TeamBID:      participatingTeams[1],
				MatchDate:    tournament.EndDate.Add(-time.Duration(3) * time.Hour * 24),
				MatchTime:    "15:00 WIB",
				Location:     "JCC Senayan - Main Stage",
				Round:        "Playoffs - Semi Final",
				Status:       "scheduled",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}

			// Semi-final 2
			semifinal2 := model.Match{
				ID:           primitive.NewObjectID(),
				TournamentID: tournament.ID,
				TeamAID:      participatingTeams[2],
				TeamBID:      participatingTeams[3],
				MatchDate:    tournament.EndDate.Add(-time.Duration(3) * time.Hour * 24),
				MatchTime:    "18:00 WIB",
				Location:     "JCC Senayan - Main Stage",
				Round:        "Playoffs - Semi Final",
				Status:       "scheduled",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}

			// Final
			final := model.Match{
				ID:           primitive.NewObjectID(),
				TournamentID: tournament.ID,
				TeamAID:      participatingTeams[0],
				TeamBID:      participatingTeams[2],
				MatchDate:    tournament.EndDate.Add(-time.Duration(1) * time.Hour * 24),
				MatchTime:    "19:00 WIB",
				Location:     "JCC Senayan - Main Stage",
				Round:        "Playoffs - Grand Final",
				Status:       "scheduled",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}

			allMatches = append(allMatches, semifinal1, semifinal2, final)
			matchCount += 3
		}
	}

	// Insert all matches
	if len(allMatches) > 0 {
		result, err := config.MatchesCollection.InsertMany(ctx, allMatches)
		if err != nil {
			log.Fatal("Failed to insert matches:", err)
		}

		fmt.Printf("✅ Successfully created %d matches across %d tournaments\n", len(result.InsertedIDs), len(tournaments))
	} else {
		fmt.Println("⚠️ No matches created")
	}

	fmt.Println("\n📊 Match Summary:")
	fmt.Printf("   - Group Stage matches: %d\n", matchCount-len(tournaments)*3)
	fmt.Printf("   - Playoff matches: %d\n", len(tournaments)*3)
	fmt.Printf("   - Total matches: %d\n", len(allMatches))

	fmt.Println("\n🎉 Match seeding completed!")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
