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

// InsertTeam creates a new team
func InsertTeam(ctx context.Context, team model.Team) (insertedID interface{}, err error) {
	// Check if team name already exists
	filter := bson.M{"team_name": team.TeamName}
	count, err := config.TeamsCollection.CountDocuments(ctx, filter)
	if err != nil {
		fmt.Printf("InsertTeam - Count Team Name: %v\n", err)
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("Team name %s sudah terdaftar", team.TeamName)
	}

	// Validate captain exists in players collection
	captainFilter := bson.M{"_id": team.CaptainID}
	captainCount, err := config.PlayersCollection.CountDocuments(ctx, captainFilter)
	if err != nil {
		fmt.Printf("InsertTeam - Check Captain: %v\n", err)
		return nil, err
	}
	if captainCount == 0 {
		return nil, fmt.Errorf("Captain dengan ID %s tidak ditemukan", team.CaptainID.Hex())
	}

	// Validate all members exist in players collection
	for _, memberID := range team.Members {
		memberFilter := bson.M{"_id": memberID}
		memberCount, err := config.PlayersCollection.CountDocuments(ctx, memberFilter)
		if err != nil {
			fmt.Printf("InsertTeam - Check Member: %v\n", err)
			return nil, err
		}
		if memberCount == 0 {
			return nil, fmt.Errorf("Member dengan ID %s tidak ditemukan", memberID.Hex())
		}
	}

	// Set timestamps
	team.CreatedAt = time.Now()
	team.UpdatedAt = time.Now()

	// Insert team
	insertResult, err := config.TeamsCollection.InsertOne(ctx, team)
	if err != nil {
		fmt.Printf("InsertTeam - Insert: %v\n", err)
		return nil, err
	}

	return insertResult.InsertedID, nil
}

// GetAllTeamsWithDetails retrieves all teams with captain details
func GetAllTeamsWithDetails(ctx context.Context) ([]model.TeamWithDetails, error) {
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "players",
				"localField":   "captain_id",
				"foreignField": "_id",
				"as":           "captain_details",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "players",
				"localField":   "members",
				"foreignField": "_id",
				"as":           "members_details",
			},
		},
		{
			"$addFields": bson.M{
				"captain_details": bson.M{
					"$arrayElemAt": []interface{}{"$captain_details", 0},
				},
			},
		},
		{
			"$project": bson.M{
				"_id":        1,
				"team_name":  1,
				"captain_id": 1,
				"members":    1,
				"logo_url":   1,
				"created_at": 1,
				"updated_at": 1,
				"captain_details": bson.M{
					"_id":         "$captain_details._id",
					"name":        "$captain_details.name",
					"ml_nickname": "$captain_details.ml_nickname",
					"ml_id":       "$captain_details.ml_id",
					"status":      "$captain_details.status",
				},
				"members_details": bson.M{
					"$map": bson.M{
						"input": "$members_details",
						"as":    "member",
						"in": bson.M{
							"_id":         "$$member._id",
							"name":        "$$member.name",
							"ml_nickname": "$$member.ml_nickname",
							"ml_id":       "$$member.ml_id",
							"status":      "$$member.status",
						},
					},
				},
			},
		},
	}

	cursor, err := config.TeamsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println("GetAllTeamsWithDetails (Aggregate):", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var teams []model.TeamWithDetails
	if err := cursor.All(ctx, &teams); err != nil {
		fmt.Println("GetAllTeamsWithDetails (Decode):", err)
		return nil, err
	}

	return teams, nil
}

// GetTeamByIDWithDetails retrieves team by ID with captain details
func GetTeamByIDWithDetails(ctx context.Context, id string) (*model.TeamWithDetails, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid team ID format")
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": objID},
		},
		{
			"$lookup": bson.M{
				"from":         "players",
				"localField":   "captain_id",
				"foreignField": "_id",
				"as":           "captain_details",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "players",
				"localField":   "members",
				"foreignField": "_id",
				"as":           "members_details",
			},
		},
		{
			"$addFields": bson.M{
				"captain_details": bson.M{
					"$arrayElemAt": []interface{}{"$captain_details", 0},
				},
			},
		},
		{
			"$project": bson.M{
				"_id":        1,
				"team_name":  1,
				"captain_id": 1,
				"members":    1,
				"logo_url":   1,
				"created_at": 1,
				"updated_at": 1,
				"captain_details": bson.M{
					"_id":         "$captain_details._id",
					"name":        "$captain_details.name",
					"ml_nickname": "$captain_details.ml_nickname",
					"ml_id":       "$captain_details.ml_id",
					"status":      "$captain_details.status",
				},
				"members_details": bson.M{
					"$map": bson.M{
						"input": "$members_details",
						"as":    "member",
						"in": bson.M{
							"_id":         "$$member._id",
							"name":        "$$member.name",
							"ml_nickname": "$$member.ml_nickname",
							"ml_id":       "$$member.ml_id",
							"status":      "$$member.status",
						},
					},
				},
			},
		},
	}

	cursor, err := config.TeamsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("terjadi kesalahan dalam mengambil data: %v", err)
	}
	defer cursor.Close(ctx)

	var teams []model.TeamWithDetails
	if err := cursor.All(ctx, &teams); err != nil {
		return nil, fmt.Errorf("terjadi kesalahan dalam decode data: %v", err)
	}

	if len(teams) == 0 {
		return nil, nil
	}

	return &teams[0], nil
}

// UpdateTeam updates team data
func UpdateTeam(ctx context.Context, id string, update model.Team) (updatedID string, err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", fmt.Errorf("invalid team ID format")
	}

	// Validate captain exists if captain_id is being updated
	if !update.CaptainID.IsZero() {
		captainFilter := bson.M{"_id": update.CaptainID}
		captainCount, err := config.PlayersCollection.CountDocuments(ctx, captainFilter)
		if err != nil {
			fmt.Printf("UpdateTeam - Check Captain: %v\n", err)
			return "", err
		}
		if captainCount == 0 {
			return "", fmt.Errorf("Captain dengan ID %s tidak ditemukan", update.CaptainID.Hex())
		}
	}

	// Validate all members exist if members are being updated
	if len(update.Members) > 0 {
		for _, memberID := range update.Members {
			memberFilter := bson.M{"_id": memberID}
			memberCount, err := config.PlayersCollection.CountDocuments(ctx, memberFilter)
			if err != nil {
				fmt.Printf("UpdateTeam - Check Member: %v\n", err)
				return "", err
			}
			if memberCount == 0 {
				return "", fmt.Errorf("Member dengan ID %s tidak ditemukan", memberID.Hex())
			}
		}
	}

	// Set updated timestamp
	update.UpdatedAt = time.Now()

	filter := bson.M{"_id": objID}
	updateData := bson.M{"$set": update}

	result, err := config.TeamsCollection.UpdateOne(ctx, filter, updateData)
	if err != nil {
		fmt.Printf("UpdateTeam: %v\n", err)
		return "", err
	}
	if result.ModifiedCount == 0 {
		return "", fmt.Errorf("tidak ada data yang diupdate untuk Team ID %s", id)
	}
	return id, nil
}

// DeleteTeam deletes team by ID
func DeleteTeam(ctx context.Context, id string) (deletedID string, err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", fmt.Errorf("invalid team ID format")
	}

	filter := bson.M{"_id": objID}
	result, err := config.TeamsCollection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Printf("DeleteTeam: %v\n", err)
		return "", err
	}
	if result.DeletedCount == 0 {
		return "", fmt.Errorf("tidak ada data yang dihapus untuk Team ID %s", id)
	}
	return id, nil
}
