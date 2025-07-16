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

// InsertPlayer creates a new player
func InsertPlayer(ctx context.Context, player model.Player) (insertedID interface{}, err error) {
	// Check if ML nickname already exists
	filter := bson.M{"ml_nickname": player.MLNickname}
	count, err := config.PlayersCollection.CountDocuments(ctx, filter)
	if err != nil {
		fmt.Printf("InsertPlayer - Count ML Nickname: %v\n", err)
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("ML Nickname %s sudah terdaftar", player.MLNickname)
	}

	// Check if ML ID already exists
	filter = bson.M{"ml_id": player.MLID}
	count, err = config.PlayersCollection.CountDocuments(ctx, filter)
	if err != nil {
		fmt.Printf("InsertPlayer - Count ML ID: %v\n", err)
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("ML ID %s sudah terdaftar", player.MLID)
	}

	// Set timestamps
	player.CreatedAt = time.Now()
	player.UpdatedAt = time.Now()

	// Insert player
	insertResult, err := config.PlayersCollection.InsertOne(ctx, player)
	if err != nil {
		fmt.Printf("InsertPlayer - Insert: %v\n", err)
		return nil, err
	}

	return insertResult.InsertedID, nil
}

// GetAllPlayers retrieves all players
func GetAllPlayers(ctx context.Context) ([]model.Player, error) {
	filter := bson.M{}

	cursor, err := config.PlayersCollection.Find(ctx, filter)
	if err != nil {
		fmt.Println("GetAllPlayers (Find):", err)
		return nil, err
	}

	var players []model.Player
	if err := cursor.All(ctx, &players); err != nil {
		fmt.Println("GetAllPlayers (Decode):", err)
		return nil, err
	}

	return players, nil
}

// GetPlayerByID retrieves player by ID
func GetPlayerByID(ctx context.Context, id string) (player *model.Player, err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid player ID format")
	}

	filter := bson.M{"_id": objID}
	err = config.PlayersCollection.FindOne(ctx, filter).Decode(&player)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("terjadi kesalahan dalam mengambil data: %v", err)
	}
	return player, nil
}

// UpdatePlayer updates player data
func UpdatePlayer(ctx context.Context, id string, update model.Player) (updatedID string, err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", fmt.Errorf("invalid player ID format")
	}

	// Set updated timestamp
	update.UpdatedAt = time.Now()

	filter := bson.M{"_id": objID}
	updateData := bson.M{"$set": update}

	result, err := config.PlayersCollection.UpdateOne(ctx, filter, updateData)
	if err != nil {
		fmt.Printf("UpdatePlayer: %v\n", err)
		return "", err
	}
	if result.ModifiedCount == 0 {
		return "", fmt.Errorf("tidak ada data yang diupdate untuk Player ID %s", id)
	}
	return id, nil
}

// DeletePlayer deletes player by ID
func DeletePlayer(ctx context.Context, id string) (deletedID string, err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", fmt.Errorf("invalid player ID format")
	}

	filter := bson.M{"_id": objID}
	result, err := config.PlayersCollection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Printf("DeletePlayer: %v\n", err)
		return "", err
	}
	if result.DeletedCount == 0 {
		return "", fmt.Errorf("tidak ada data yang dihapus untuk Player ID %s", id)
	}
	return id, nil
}
