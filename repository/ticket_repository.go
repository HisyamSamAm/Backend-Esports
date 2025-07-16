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

// CreateTicket creates a new ticket type
func CreateTicket(ctx context.Context, ticket model.Ticket) (insertedID interface{}, err error) {
	// Validate tournament exists
	tournamentFilter := bson.M{"_id": ticket.TournamentID}
	tournamentCount, err := config.TournamentsCollection.CountDocuments(ctx, tournamentFilter)
	if err != nil {
		fmt.Printf("CreateTicket - Check Tournament: %v\n", err)
		return nil, err
	}
	if tournamentCount == 0 {
		return nil, fmt.Errorf("Tournament dengan ID %s tidak ditemukan", ticket.TournamentID.Hex())
	}

	// Set timestamps
	ticket.CreatedAt = time.Now()
	ticket.UpdatedAt = time.Now()

	// Insert ticket
	insertResult, err := config.TicketsCollection.InsertOne(ctx, ticket)
	if err != nil {
		fmt.Printf("CreateTicket - Insert: %v\n", err)
		return nil, err
	}

	return insertResult.InsertedID, nil
}

// GetAllTickets retrieves all tickets with optional tournament filter
func GetAllTickets(ctx context.Context, tournamentID string) ([]model.Ticket, error) {
	filter := bson.M{}

	// Add tournament filter if provided
	if tournamentID != "" {
		objID, err := primitive.ObjectIDFromHex(tournamentID)
		if err != nil {
			return nil, fmt.Errorf("invalid tournament ID format")
		}
		filter["tournament_id"] = objID
	}

	cursor, err := config.TicketsCollection.Find(ctx, filter)
	if err != nil {
		fmt.Println("GetAllTickets (Find):", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var tickets []model.Ticket
	if err := cursor.All(ctx, &tickets); err != nil {
		fmt.Println("GetAllTickets (Decode):", err)
		return nil, err
	}

	return tickets, nil
}

// GetTicketByID retrieves ticket by ID
func GetTicketByID(ctx context.Context, id string) (*model.Ticket, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ticket ID format")
	}

	var ticket model.Ticket
	filter := bson.M{"_id": objID}
	err = config.TicketsCollection.FindOne(ctx, filter).Decode(&ticket)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("terjadi kesalahan dalam mengambil data: %v", err)
	}
	return &ticket, nil
}

// GetTicketWithTournamentByID retrieves ticket by ID with tournament details populated
func GetTicketWithTournamentByID(ctx context.Context, id string) (*model.TicketWithTournament, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ticket ID format")
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": objID},
		},
		{
			"$lookup": bson.M{
				"from":         "tournaments",
				"localField":   "tournament_id",
				"foreignField": "_id",
				"as":           "tournament_details",
			},
		},
		{
			"$project": bson.M{
				"_id":                1,
				"tournament_id":      1,
				"price":              1,
				"quantity_available": 1,
				"description":        1,
				"created_at":         1,
				"updated_at":         1,
				"tournament": bson.M{
					"$arrayElemAt": []interface{}{
						bson.M{
							"$map": bson.M{
								"input": "$tournament_details",
								"as":    "tournament",
								"in": bson.M{
									"_id":    "$$tournament._id",
									"name":   "$$tournament.name",
									"status": "$$tournament.status",
								},
							},
						}, 0,
					},
				},
			},
		},
	}

	cursor, err := config.TicketsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("terjadi kesalahan dalam mengambil data: %v", err)
	}
	defer cursor.Close(ctx)

	var tickets []model.TicketWithTournament
	if err := cursor.All(ctx, &tickets); err != nil {
		return nil, fmt.Errorf("terjadi kesalahan dalam decode data: %v", err)
	}

	if len(tickets) == 0 {
		return nil, nil
	}

	return &tickets[0], nil
}

// GetAllTicketsWithTournament retrieves all tickets with tournament details populated
func GetAllTicketsWithTournament(ctx context.Context, tournamentID string) ([]model.TicketWithTournament, error) {
	filter := bson.M{}
	if tournamentID != "" {
		tournamentObjectID, err := primitive.ObjectIDFromHex(tournamentID)
		if err != nil {
			return nil, fmt.Errorf("invalid tournament ID format: %v", err)
		}
		filter["tournament_id"] = tournamentObjectID
	}

	// Aggregation pipeline to join with tournaments collection
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{
			{Key: "$lookup", Value: bson.M{
				"from":         "tournaments",
				"localField":   "tournament_id",
				"foreignField": "_id",
				"as":           "tournament",
			}},
		},
		{
			{Key: "$addFields", Value: bson.M{
				"tournament": bson.M{
					"$cond": bson.M{
						"if": bson.M{"$gt": bson.A{bson.M{"$size": "$tournament"}, 0}},
						"then": bson.M{
							"_id":        bson.M{"$arrayElemAt": bson.A{"$tournament._id", 0}},
							"name":       bson.M{"$arrayElemAt": bson.A{"$tournament.name", 0}},
							"game":       bson.M{"$arrayElemAt": bson.A{"$tournament.game", 0}},
							"type":       bson.M{"$arrayElemAt": bson.A{"$tournament.type", 0}},
							"start_date": bson.M{"$arrayElemAt": bson.A{"$tournament.start_date", 0}},
						},
						"else": nil,
					},
				},
			}},
		},
	}

	cursor, err := config.TicketsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("terjadi kesalahan dalam mengambil data: %v", err)
	}
	defer cursor.Close(ctx)

	var tickets []model.TicketWithTournament
	if err := cursor.All(ctx, &tickets); err != nil {
		return nil, fmt.Errorf("terjadi kesalahan dalam decode data: %v", err)
	}

	return tickets, nil
}

// UpdateTicket updates ticket data
func UpdateTicket(ctx context.Context, id string, update model.Ticket) (updatedID string, err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", fmt.Errorf("invalid ticket ID format")
	}

	// Validate tournament exists if tournament_id is being updated
	if !update.TournamentID.IsZero() {
		tournamentFilter := bson.M{"_id": update.TournamentID}
		tournamentCount, err := config.TournamentsCollection.CountDocuments(ctx, tournamentFilter)
		if err != nil {
			fmt.Printf("UpdateTicket - Check Tournament: %v\n", err)
			return "", err
		}
		if tournamentCount == 0 {
			return "", fmt.Errorf("Tournament dengan ID %s tidak ditemukan", update.TournamentID.Hex())
		}
	}

	// Set updated timestamp
	update.UpdatedAt = time.Now()

	filter := bson.M{"_id": objID}
	updateData := bson.M{"$set": update}

	result, err := config.TicketsCollection.UpdateOne(ctx, filter, updateData)
	if err != nil {
		fmt.Printf("UpdateTicket: %v\n", err)
		return "", err
	}
	if result.ModifiedCount == 0 {
		return "", fmt.Errorf("tidak ada data yang diupdate untuk Ticket ID %s", id)
	}
	return id, nil
}

// DeleteTicket deletes ticket by ID
func DeleteTicket(ctx context.Context, id string) (deletedID string, err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", fmt.Errorf("invalid ticket ID format")
	}

	filter := bson.M{"_id": objID}
	result, err := config.TicketsCollection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Printf("DeleteTicket: %v\n", err)
		return "", err
	}
	if result.DeletedCount == 0 {
		return "", fmt.Errorf("tidak ada data yang dihapus untuk Ticket ID %s", id)
	}
	return id, nil
}

// GetTicketsByTournamentID retrieves all tickets for a specific tournament
func GetTicketsByTournamentID(ctx context.Context, tournamentID string) ([]model.Ticket, error) {
	objID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return nil, fmt.Errorf("invalid tournament ID format")
	}

	filter := bson.M{"tournament_id": objID}
	cursor, err := config.TicketsCollection.Find(ctx, filter)
	if err != nil {
		fmt.Println("GetTicketsByTournamentID (Find):", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var tickets []model.Ticket
	if err := cursor.All(ctx, &tickets); err != nil {
		fmt.Println("GetTicketsByTournamentID (Decode):", err)
		return nil, err
	}

	return tickets, nil
}

// GetTicketsByTournament gets all tickets for a specific tournament
func GetTicketsByTournament(tournamentID primitive.ObjectID) ([]model.Ticket, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"tournament_id": tournamentID}
	cursor, err := config.TicketsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tickets []model.Ticket
	if err = cursor.All(ctx, &tickets); err != nil {
		return nil, err
	}

	return tickets, nil
}

// GetTicketsByUser gets all tickets owned by a specific user
func GetTicketsByUser(userID primitive.ObjectID) ([]model.Ticket, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	cursor, err := config.TicketsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tickets []model.Ticket
	if err = cursor.All(ctx, &tickets); err != nil {
		return nil, err
	}

	return tickets, nil
}
