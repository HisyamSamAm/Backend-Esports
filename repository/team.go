package controller

import (
	"EMBECK/config"
	"EMBECK/model"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func GetAllTeams(ctx context.Context) ([]model.Team, error) {
	collection := config.DB.Collection("team")
	filter := bson.M{}

	cursor, err:= collection.Find(ctx, filter)
	if err != nil {
		fmt.Println("GetAllTeams (find):", err)
		return nil, err
	}

	var data []model.Team
	if err := cursor.All(ctx, &data); err != nil {
		fmt.Println("GetAllTeams (Decode):", err)
		return nil, err
	}

	return data, nil
}

func GetTeamByID(ctx context.Context, id string) (model.Team, error) {
	collection := config.DB.Collection("team")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Invalid ID format:", err)
		return model.Team{}, err
	}

	var team model.Team 
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&team)
	if err != nil {
		fmt.Println("GetteamByID (Find):", err)
		return model.Team{}, err
	}

	return team, nil
}

func CreateTeam(ctx context.Context, team model.Team) (interface{}, error) {
	collection := config.DB.Collection("team")

	result, err := collection.InsertOne(ctx, team)
	if err != nil {
		fmt.Println("CreateTeam:", err)
		return nil, err
	}

	return result.InsertedID, nil
}

func UpdateTeam(ctx context.Context, id string, team model.Team) error {
	collection := config.DB.Collection("team")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Invalid ID:", err)
		return err
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": team})
	if err != nil {
		fmt.Println("UpdateTeam:", err)
		return err
	}

	return nil
}

func DeleteTeam(ctx context.Context, id string) error {
	collection := config.DB.Collection("team")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Invalid ID:", err)
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		fmt.Println("DeleteTeam:", err)
		return err
	}

	return nil
}

