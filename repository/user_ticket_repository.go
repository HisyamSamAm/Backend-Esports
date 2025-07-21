package repository

import (
	"context"
	"embeck/config"
	"embeck/model"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PurchaseTicket creates a new ticket record for a user and match.
func PurchaseTicket(ctx context.Context, userID, matchID primitive.ObjectID) (*model.UserTicket, error) {
	// 1. Validate if the match exists
	matchCount, err := config.MatchesCollection.CountDocuments(ctx, bson.M{"_id": matchID})
	if err != nil {
		return nil, fmt.Errorf("error validating match: %w", err)
	}
	if matchCount == 0 {
		return nil, fmt.Errorf("match not found")
	}

	// 2. (Optional) Check if user has already bought a ticket for this match
	// This prevents a user from buying multiple tickets for the same match.
	// You can remove this check if users are allowed to buy more than one.
	filter := bson.M{"user_id": userID, "match_id": matchID}
	existingTicketCount, err := config.UserTicketsCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error checking existing tickets: %w", err)
	}
	if existingTicketCount > 0 {
		return nil, fmt.Errorf("ticket for this match already purchased")
	}

	// 3. Create the new user ticket
	newUserTicket := model.UserTicket{
		UserID:       userID,
		MatchID:      matchID,
		PurchaseDate: time.Now(),
		Status:       "valid", // Default status upon purchase
	}

	result, err := config.UserTicketsCollection.InsertOne(ctx, newUserTicket)
	if err != nil {
		return nil, fmt.Errorf("failed to insert ticket: %w", err)
	}

	newUserTicket.ID = result.InsertedID.(primitive.ObjectID)
	return &newUserTicket, nil
}

// GetTicketsByUserID retrieves all tickets for a specific user with populated match details.
func GetTicketsByUserID(ctx context.Context, userID primitive.ObjectID) ([]model.UserTicketResponse, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{"user_id": userID},
		},
		{
			"$lookup": bson.M{
				"from":         "matches",
				"localField":   "match_id",
				"foreignField": "_id",
				"as":           "match_details_full",
			},
		},
		{
			"$unwind": "$match_details_full",
		},
		{
			"$lookup": bson.M{
				"from":         "teams",
				"localField":   "match_details_full.team_a_id",
				"foreignField": "_id",
				"as":           "team_a_info",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "teams",
				"localField":   "match_details_full.team_b_id",
				"foreignField": "_id",
				"as":           "team_b_info",
			},
		},
		{
			"$project": bson.M{
				"_id":           1,
				"user_id":       1,
				"match_id":      1,
				"purchase_date": 1,
				"status":        1,
				"match_details": bson.M{
					"id":         "$match_details_full._id",
					"match_date": "$match_details_full.match_date",
					"match_time": "$match_details_full.match_time",
					"round":      "$match_details_full.round",
					"status":     "$match_details_full.status",
					"team_a":     bson.M{"$arrayElemAt": []interface{}{"$team_a_info", 0}},
					"team_b":     bson.M{"$arrayElemAt": []interface{}{"$team_b_info", 0}},
				},
			},
		},
		{
			// Cleanup team fields to match TeamBasicInfo
			"$project": bson.M{
				"_id":                            1,
				"user_id":                        1,
				"match_id":                       1,
				"purchase_date":                  1,
				"status":                         1,
				"match_details.id":               1,
				"match_details.match_date":       1,
				"match_details.match_time":       1,
				"match_details.round":            1,
				"match_details.status":           1,
				"match_details.team_a._id":       "$match_details.team_a._id",
				"match_details.team_a.team_name": "$match_details.team_a.team_name",
				"match_details.team_a.logo_url":  "$match_details.team_a.logo_url",
				"match_details.team_b._id":       "$match_details.team_b._id",
				"match_details.team_b.team_name": "$match_details.team_b.team_name",
				"match_details.team_b.logo_url":  "$match_details.team_b.logo_url",
			},
		},
	}

	cursor, err := config.UserTicketsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregation failed: %w", err)
	}
	defer cursor.Close(ctx)

	var tickets []model.UserTicketResponse
	if err := cursor.All(ctx, &tickets); err != nil {
		return nil, fmt.Errorf("failed to decode tickets: %w", err)
	}

	return tickets, nil
}
