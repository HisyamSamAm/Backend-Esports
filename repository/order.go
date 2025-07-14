package controller

import (
	"EMBECK/config"
	"EMBECK/model"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func GetAllOrders(ctx context.Context) ([]model.Order, error) {
	collection := config.DB.Collection("orders")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("GetAllOrders (Find):", err)
		return nil, err
	}

	var orders []model.Order
	if err := cursor.All(ctx, &orders); err != nil {
		fmt.Println("GetAllOrders (Decode):", err)
		return nil, err
	}

	return orders, nil
}
func GetOrderByID(ctx context.Context, id string) (model.Order, error) {
	collection := config.DB.Collection("orders")

	filter := bson.M{"id": id} // cari berdasarkan field "id"
	var order model.Order
	err := collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		fmt.Println("GetOrderByID:", err)
		return model.Order{}, err
	}

	return order, nil
}
func CreateOrder(ctx context.Context, order model.Order) error {
	collection := config.DB.Collection("orders")

	_, err := collection.InsertOne(ctx, order)
	if err != nil {
		fmt.Println("CreateOrder:", err)
		return err
	}
	return nil
}
func UpdateOrder(ctx context.Context, id string, updatedData model.Order) error {
	collection := config.DB.Collection("orders")

	filter := bson.M{"id": id}
	update := bson.M{
		"$set": updatedData,
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("UpdateOrder:", err)
		return err
	}
	return nil
}
func DeleteOrder(ctx context.Context, id string) error {
	collection := config.DB.Collection("orders")

	filter := bson.M{"id": id}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println("DeleteOrder:", err)
		return err
	}
	return nil
}