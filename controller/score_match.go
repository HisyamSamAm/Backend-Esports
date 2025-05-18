package controller

import (
	"EMBECK/config"
	"EMBECK/model"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func GetAllScoreMatches(ctx context.Context) ([]model.Score, error) {
	collection := config.DB.Collection("score")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("GetAllScoreMatches (Find):", err)
		return nil, err
	}

	var matches []model.Score
	if err := cursor.All(ctx, &matches); err != nil {
		fmt.Println("GetAllScoreMatches (Decode):", err)
		return nil, err
	}

	return matches, nil
}


func GetAllScoreMatchesByID(ctx context.Context, id string) (model.Score, error) {
	collection := config.DB.Collection("score")

	filter := bson.M{"id": id} // cari berdasarkan field "id"
	var scoreMatch model.Score
	err := collection.FindOne(ctx, filter).Decode(&scoreMatch)
	if err != nil {
		fmt.Println("GetAllScoreMatchesByID:", err)
		return model.Score{}, err
	}

	return scoreMatch, nil
}

func CreateScoreMatch(ctx context.Context, scoreMatch model.Score) error {
	collection := config.DB.Collection("score")

	_, err := collection.InsertOne(ctx, scoreMatch)
	if err != nil {
		fmt.Println("CreateScoreMatch:", err)
		return err
	}
	return nil
}

func UpdateScoreMatch(ctx context.Context, id string, updatedData model.Score) error {
	collection := config.DB.Collection("score")

	filter := bson.M{"id": id}
	update := bson.M{
		"$set": updatedData,
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("UpdateScoreMatch:", err)
		return err
	}
	return nil
}

func DeleteScoreMatch(ctx context.Context, id string) error {
	collection := config.DB.Collection("score")

	filter := bson.M{"id": id}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println("DeleteScoreMatch:", err)
		return err
	}
	return nil
}