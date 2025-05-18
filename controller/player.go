package controller

import (
	"EMBECK/config"
	"EMBECK/model"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func GetAllPlayers(ctx context.Context) ([]model.Player, error) {
	collection := config.DB.Collection("player")

	cursor, err := collection.Find(ctx, bson.M{})
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


func GetPlayerByID(ctx context.Context, id string) (model.Player, error) {
	collection := config.DB.Collection("player")

	filter := bson.M{"id": id} // cari berdasarkan field "id"
	var player model.Player
	err := collection.FindOne(ctx, filter).Decode(&player)
	if err != nil {
		fmt.Println("GetPlayerByID:", err)
		return model.Player{}, err
	}

	return player, nil
}


func CreatePlayer(ctx context.Context, player model.Player) error {
	collection := config.DB.Collection("player")

	_, err := collection.InsertOne(ctx, player)
	if err != nil {
		fmt.Println("CreatePlayer:", err)
		return err
	}
	return nil
}

func UpdatePlayer(ctx context.Context, id string, updatedData model.Player) error {
	collection := config.DB.Collection("players")

	filter := bson.M{"id": id}

	// Jangan update field ID biar gak ketiban null
	update := bson.M{
		"$set": bson.M{
			"name":         updatedData.Name,
			"in_game_name": updatedData.InGameName,
			"role":         updatedData.Role,
			"team_id":      updatedData.TeamID,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("UpdatePlayer:", err)
		return err
	}
	return nil
}



func DeletePlayer(ctx context.Context, id string) error {
	collection := config.DB.Collection("player")

	filter := bson.M{"id": id}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println("DeletePlayer:", err)
		return err
	}
	return nil
}

