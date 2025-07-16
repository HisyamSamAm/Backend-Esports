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

	fmt.Println("🎫 Creating sample tickets for tournaments...")

	ctx := context.Background()

	// Get existing tournaments to create tickets for
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

	// Clear existing tickets
	_, err = config.TicketsCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to clear existing tickets:", err)
	}

	// Create sample tickets for each tournament
	var allTickets []interface{}
	ticketCount := 0

	for _, tournament := range tournaments {
		// Create different ticket types for each tournament
		ticketTypes := []struct {
			Price       int
			Quantity    int
			Description string
		}{
			{50000, 1000, "Tiket Reguler - Tribun A (Standing)"},
			{75000, 800, "Tiket Reguler - Tribun B (Sitting)"},
			{100000, 500, "Tiket Premium - Tribun VIP (Sitting dengan fasilitas)"},
			{150000, 200, "Tiket VVIP - Akses backstage dan meet & greet"},
			{25000, 1500, "Tiket Early Bird - Akses hari pertama saja"},
			{300000, 50, "Tiket Corporate - Package untuk 5 orang"},
		}

		for i, ticketType := range ticketTypes {
			ticket := model.Ticket{
				ID:                primitive.NewObjectID(),
				TournamentID:      tournament.ID,
				Price:             ticketType.Price,
				QuantityAvailable: ticketType.Quantity,
				Description:       ticketType.Description,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}

			allTickets = append(allTickets, ticket)
			ticketCount++

			fmt.Printf("  📋 %s - %s: Rp %d (%d tersedia)\n",
				tournament.Name,
				ticketType.Description,
				ticketType.Price,
				ticketType.Quantity)

			// Limit to avoid too many tickets
			if i >= 3 { // Only create first 4 ticket types per tournament
				break
			}
		}
	}

	// Insert all tickets
	if len(allTickets) > 0 {
		result, err := config.TicketsCollection.InsertMany(ctx, allTickets)
		if err != nil {
			log.Fatal("Failed to insert tickets:", err)
		}

		fmt.Printf("\n✅ Successfully created %d ticket types across %d tournaments\n", len(result.InsertedIDs), len(tournaments))
	} else {
		fmt.Println("⚠️ No tickets created")
	}

	// Calculate total capacity and revenue potential
	totalCapacity := 0
	totalRevenuePotential := 0

	for _, ticketInterface := range allTickets {
		ticket := ticketInterface.(model.Ticket)
		totalCapacity += ticket.QuantityAvailable
		totalRevenuePotential += (ticket.Price * ticket.QuantityAvailable)
	}

	fmt.Println("\n📊 Ticket Summary:")
	fmt.Printf("   - Total ticket types: %d\n", len(allTickets))
	fmt.Printf("   - Total capacity: %d tiket\n", totalCapacity)
	fmt.Printf("   - Revenue potential: Rp %s\n", formatRupiah(totalRevenuePotential))
	fmt.Printf("   - Average ticket price: Rp %s\n", formatRupiah(totalRevenuePotential/totalCapacity))

	fmt.Println("\n🎫 Ticket seeding completed!")
}

func formatRupiah(amount int) string {
	if amount >= 1000000000 {
		return fmt.Sprintf("%.1f miliar", float64(amount)/1000000000)
	} else if amount >= 1000000 {
		return fmt.Sprintf("%.1f juta", float64(amount)/1000000)
	} else if amount >= 1000 {
		return fmt.Sprintf("%.1f ribu", float64(amount)/1000)
	}
	return fmt.Sprintf("%d", amount)
}
