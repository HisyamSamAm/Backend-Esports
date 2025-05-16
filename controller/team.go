package controller

import (
	"EMBECK/config"
	"EMBECK/model"
	"context"
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

func GetAllTeams(ctx context.Context) ([]model.Team, error) {
	collection := config.DB.Collection("teams")
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

