package controller

import (
	"EMBECK/config"
	"EMBECK/model"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func GetAllUsers(ctx context.Context) ([]model.User, error) {
	collection := config.DB.Collection("users")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("GetAllUsers (Find):", err)
		return nil, err
	}

	var users []model.User
	if err := cursor.All(ctx, &users); err != nil {
		fmt.Println("GetAllUsers (Decode):", err)
		return nil, err
	}

	return users, nil
}

func GetUserByID(ctx context.Context, id string) (model.User, error) {
	collection := config.DB.Collection("users")

	filter := bson.M{"id": id} // cari berdasarkan field "id"
	var user model.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		fmt.Println("GetUserByID:", err)
		return model.User{}, err
	}

	return user, nil
}

func CreateUser(ctx context.Context, user model.User) error {
	collection := config.DB.Collection("users")

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		fmt.Println("CreateUser:", err)
		return err
	}
	return nil
}

func UpdateUser(ctx context.Context, id string, updatedData model.User) error {
	collection := config.DB.Collection("users")

	filter := bson.M{"id": id}
	update := bson.M{
		"$set": updatedData,
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("UpdateUser:", err)
		return err
	}
	return nil
}

func DeleteUser(ctx context.Context, id string) error {
	collection := config.DB.Collection("users")

	filter := bson.M{"id": id}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println("DeleteUser:", err)
		return err
	}
	return nil
}

func AuthenticateUser(ctx context.Context, username, password string) (model.User, error) {
	collection := config.DB.Collection("users")

	filter := bson.M{"username": username, "password": password}
	var user model.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		fmt.Println("AuthenticateUser:", err)
		return model.User{}, err
	}

	return user, nil
}