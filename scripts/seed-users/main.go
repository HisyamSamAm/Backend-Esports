package main

import (
	"context"
	"embeck/config"
	"embeck/model"
	"embeck/pkg/password"
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

	fmt.Println("👤 Creating sample users...")

	ctx := context.Background()

	// Clear existing users
	_, err := config.UsersCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to clear existing users:", err)
	}
	fmt.Println("🧹 Cleared existing users")

	// Sample users to create
	usersData := []struct {
		Username string
		Email    string
		Password string
		Role     string
	}{
		{
			Username: "admin",
			Email:    "admin@embeck.com",
			Password: "admin123",
			Role:     "admin",
		},
		{
			Username: "moderator",
			Email:    "moderator@embeck.com",
			Password: "moderator123",
			Role:     "admin",
		},
		{
			Username: "user1",
			Email:    "user1@example.com",
			Password: "user123",
			Role:     "user",
		},
		{
			Username: "user2",
			Email:    "user2@example.com",
			Password: "user123",
			Role:     "user",
		},
		{
			Username: "spectator",
			Email:    "spectator@example.com",
			Password: "spectator123",
			Role:     "user",
		},
	}

	// Create users
	var allUsers []interface{}
	userCount := 0

	for _, userData := range usersData {
		// Hash password
		hashedPassword, err := password.HashPassword(userData.Password)
		if err != nil {
			log.Printf("Failed to hash password for %s: %v", userData.Username, err)
			continue
		}

		user := model.User{
			ID:        primitive.NewObjectID(),
			Username:  userData.Username,
			Email:     userData.Email,
			Password:  hashedPassword,
			Role:      userData.Role,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		allUsers = append(allUsers, user)
		userCount++

		fmt.Printf("  👤 %s (%s) - %s\n", userData.Username, userData.Role, userData.Email)
	}

	// Insert all users
	if len(allUsers) > 0 {
		_, err := config.UsersCollection.InsertMany(ctx, allUsers)
		if err != nil {
			log.Fatal("Failed to insert users:", err)
		}
		fmt.Printf("✅ Successfully created %d users\n", len(allUsers))
	} else {
		fmt.Println("❌ No users to insert")
	}

	fmt.Println("\n📊 User Summary:")
	adminCount := 0
	userCountByRole := 0
	for _, userData := range usersData {
		if userData.Role == "admin" {
			adminCount++
		} else {
			userCountByRole++
		}
	}
	fmt.Printf("   - Admin users: %d\n", adminCount)
	fmt.Printf("   - Regular users: %d\n", userCountByRole)
	fmt.Printf("   - Total users: %d\n", len(allUsers))

	fmt.Println("\n🔐 Default Login Credentials:")
	fmt.Println("   Admin:")
	fmt.Println("     Email: admin@embeck.com")
	fmt.Println("     Password: admin123")
	fmt.Println("   User:")
	fmt.Println("     Email: user1@example.com")
	fmt.Println("     Password: user123")

	fmt.Println("\n👤 User seeding completed!")
}
