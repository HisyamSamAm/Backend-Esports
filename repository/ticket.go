package controller

import (
	"EMBECK/config"
	"EMBECK/model"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func GetAllTickets(ctx context.Context) ([]model.Ticket, error) {
	collection := config.DB.Collection("ticket")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("GetAllTickets (Find):", err)
		return nil, err
	}

	var tickets []model.Ticket
	if err := cursor.All(ctx, &tickets); err != nil {
		fmt.Println("GetAllTickets (Decode):", err)
		return nil, err
	}

	return tickets, nil
}
func GetTicketByID(ctx context.Context, id string) (model.Ticket, error) {
	collection := config.DB.Collection("ticket")

	filter := bson.M{"id": id} // cari berdasarkan field "id"
	var ticket model.Ticket
	err := collection.FindOne(ctx, filter).Decode(&ticket)
	if err != nil {
		fmt.Println("GetTicketByID:", err)
		return model.Ticket{}, err
	}

	return ticket, nil
}

func CreateTicket(ctx context.Context, ticket model.Ticket) error {
	collection := config.DB.Collection("ticket")

	_, err := collection.InsertOne(ctx, ticket)
	if err != nil {
		fmt.Println("CreateTicket:", err)
		return err
	}
	return nil
}

func UpdateTicket(ctx context.Context, id string, updatedData model.Ticket) error {
	collection := config.DB.Collection("ticket")

	filter := bson.M{"id": id}
	update := bson.M{
		"$set": updatedData,
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("UpdateTicket:", err)
		return err
	}
	return nil
}

func DeleteTicket(ctx context.Context, id string) error {
	collection := config.DB.Collection("ticket")

	filter := bson.M{"id": id}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println("DeleteTicket:", err)
		return err
	}
	return nil
}