package repository

import (
	"context"
	"embeck/config"
	"embeck/model"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateUser creates a new user with dual sync support
func CreateUser(ctx context.Context, user model.User) (insertedID interface{}, err error) {
	// Check if username already exists
	usernameFilter := bson.M{"username": user.Username}
	usernameCount, err := config.UsersCollection.CountDocuments(ctx, usernameFilter)
	if err != nil {
		fmt.Printf("CreateUser - Check Username: %v\n", err)
		return nil, err
	}
	if usernameCount > 0 {
		return nil, fmt.Errorf("Username %s sudah terdaftar", user.Username)
	}

	// Check if email already exists
	emailFilter := bson.M{"email": user.Email}
	emailCount, err := config.UsersCollection.CountDocuments(ctx, emailFilter)
	if err != nil {
		fmt.Printf("CreateUser - Check Email: %v\n", err)
		return nil, err
	}
	if emailCount > 0 {
		return nil, fmt.Errorf("Email %s sudah terdaftar", user.Email)
	}

	// Set timestamps
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Insert user with dual sync capability
	insertResult, err := config.UsersCollection.InsertOne(ctx, user)
	if err != nil {
		fmt.Printf("CreateUser - Insert: %v\n", err)
		return nil, err
	}

	return insertResult.InsertedID, nil
}

// GetUserByEmail retrieves user by email
func GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	filter := bson.M{"email": email}
	err := config.UsersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("terjadi kesalahan dalam mengambil data user: %v", err)
	}
	return &user, nil
}

// GetUserByID retrieves user by ID
func GetUserByID(ctx context.Context, id string) (*model.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format")
	}

	var user model.User
	filter := bson.M{"_id": objID}
	err = config.UsersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("terjadi kesalahan dalam mengambil data user: %v", err)
	}
	return &user, nil
}

// GetUserByUsername retrieves user by username
func GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	filter := bson.M{"username": username}
	err := config.UsersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("terjadi kesalahan dalam mengambil data user: %v", err)
	}
	return &user, nil
}

// GetAllUsers retrieves all users (admin only)
func GetAllUsers(ctx context.Context) ([]model.UserProfile, error) {
	cursor, err := config.UsersCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("GetAllUsers (Find):", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []model.User
	if err := cursor.All(ctx, &users); err != nil {
		fmt.Println("GetAllUsers (Decode):", err)
		return nil, err
	}

	// Convert to UserProfile (without password)
	var userProfiles []model.UserProfile
	for _, user := range users {
		profile := model.UserProfile{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		userProfiles = append(userProfiles, profile)
	}

	return userProfiles, nil
}

// UpdateUser updates user data
func UpdateUser(ctx context.Context, id string, update model.User) (updatedID string, err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", fmt.Errorf("invalid user ID format")
	}

	// Check if username already exists (excluding current user)
	if update.Username != "" {
		usernameFilter := bson.M{
			"username": update.Username,
			"_id":      bson.M{"$ne": objID},
		}
		usernameCount, err := config.UsersCollection.CountDocuments(ctx, usernameFilter)
		if err != nil {
			fmt.Printf("UpdateUser - Check Username: %v\n", err)
			return "", err
		}
		if usernameCount > 0 {
			return "", fmt.Errorf("Username %s sudah digunakan user lain", update.Username)
		}
	}

	// Check if email already exists (excluding current user)
	if update.Email != "" {
		emailFilter := bson.M{
			"email": update.Email,
			"_id":   bson.M{"$ne": objID},
		}
		emailCount, err := config.UsersCollection.CountDocuments(ctx, emailFilter)
		if err != nil {
			fmt.Printf("UpdateUser - Check Email: %v\n", err)
			return "", err
		}
		if emailCount > 0 {
			return "", fmt.Errorf("Email %s sudah digunakan user lain", update.Email)
		}
	}

	// Set updated timestamp
	update.UpdatedAt = time.Now()

	filter := bson.M{"_id": objID}
	updateData := bson.M{"$set": update}

	result, err := config.UsersCollection.UpdateOne(ctx, filter, updateData)
	if err != nil {
		fmt.Printf("UpdateUser: %v\n", err)
		return "", err
	}
	if result.ModifiedCount == 0 {
		return "", fmt.Errorf("tidak ada data yang diupdate untuk User ID %s", id)
	}
	return id, nil
}

// DeleteUser deletes user by ID
func DeleteUser(ctx context.Context, id string) (deletedID string, err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", fmt.Errorf("invalid user ID format")
	}

	filter := bson.M{"_id": objID}
	result, err := config.UsersCollection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Printf("DeleteUser: %v\n", err)
		return "", err
	}
	if result.DeletedCount == 0 {
		return "", fmt.Errorf("tidak ada data yang dihapus untuk User ID %s", id)
	}
	return id, nil
}
