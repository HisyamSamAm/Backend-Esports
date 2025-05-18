package controller

import (
	"EMBECK/config"
	"EMBECK/model"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func GetAllTournaments(ctx context.Context) ([]model.Tournament, error) {
	collection := config.DB.Collection("tournaments")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("GetAllTournaments (Find):", err)
		return nil, err
	}

	var tournaments []model.Tournament
	if err := cursor.All(ctx, &tournaments); err != nil {
		fmt.Println("GetAllTournaments (Decode):", err)
		return nil, err
	}

	return tournaments, nil
}

func GetTournamentByID(ctx context.Context, id string) (model.Tournament, error) {
	collection := config.DB.Collection("tournaments")

	filter := bson.M{"id": id} // cari berdasarkan field "id"
	var tournament model.Tournament
	err := collection.FindOne(ctx, filter).Decode(&tournament)
	if err != nil {
		fmt.Println("GetTournamentByID:", err)
		return model.Tournament{}, err
	}

	return tournament, nil
}

func CreateTournament(ctx context.Context, tournament model.Tournament) error {
	collection := config.DB.Collection("tournaments")

	_, err := collection.InsertOne(ctx, tournament)
	if err != nil {
		fmt.Println("CreateTournament:", err)
		return err
	}
	return nil
}


func UpdateTournament(ctx context.Context, id string, updatedData model.Tournament) error {
	collection := config.DB.Collection("tournaments")

	filter := bson.M{"id": id}

	update := bson.M{
		"$set": updatedData,
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("UpdateTournament:", err)
		return err
	}
	return nil
}

func DeleteTournament(ctx context.Context, id string) error {
	collection := config.DB.Collection("tournaments")

	filter := bson.M{"id": id}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println("DeleteTournament:", err)
		return err
	}	
	return nil
}