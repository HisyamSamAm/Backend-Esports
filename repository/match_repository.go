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

// CreateMatch creates a new match
func CreateMatch(ctx context.Context, match model.Match) (insertedID interface{}, err error) {
	// Validate tournament exists
	tournamentFilter := bson.M{"_id": match.TournamentID}
	tournamentCount, err := config.TournamentsCollection.CountDocuments(ctx, tournamentFilter)
	if err != nil {
		fmt.Printf("CreateMatch - Check Tournament: %v\n", err)
		return nil, err
	}
	if tournamentCount == 0 {
		return nil, fmt.Errorf("Tournament dengan ID %s tidak ditemukan", match.TournamentID.Hex())
	}

	// Validate team A exists
	teamAFilter := bson.M{"_id": match.TeamAID}
	teamACount, err := config.TeamsCollection.CountDocuments(ctx, teamAFilter)
	if err != nil {
		fmt.Printf("CreateMatch - Check Team A: %v\n", err)
		return nil, err
	}
	if teamACount == 0 {
		return nil, fmt.Errorf("Team A dengan ID %s tidak ditemukan", match.TeamAID.Hex())
	}

	// Validate team B exists
	teamBFilter := bson.M{"_id": match.TeamBID}
	teamBCount, err := config.TeamsCollection.CountDocuments(ctx, teamBFilter)
	if err != nil {
		fmt.Printf("CreateMatch - Check Team B: %v\n", err)
		return nil, err
	}
	if teamBCount == 0 {
		return nil, fmt.Errorf("Team B dengan ID %s tidak ditemukan", match.TeamBID.Hex())
	}

	// Validate teams are different
	if match.TeamAID == match.TeamBID {
		return nil, fmt.Errorf("Team A dan Team B harus berbeda")
	}

	// Set timestamps
	match.CreatedAt = time.Now()
	match.UpdatedAt = time.Now()

	// Insert match
	insertResult, err := config.MatchesCollection.InsertOne(ctx, match)
	if err != nil {
		fmt.Printf("CreateMatch - Insert: %v\n", err)
		return nil, err
	}

	return insertResult.InsertedID, nil
}

// GetAllMatches retrieves all matches with populated team details
func GetAllMatches(ctx context.Context, tournamentID string) ([]model.MatchWithDetails, error) {
	pipeline := []bson.M{}

	// Add tournament filter if provided
	if tournamentID != "" && tournamentID != "all" {
		objID, err := primitive.ObjectIDFromHex(tournamentID)
		if err != nil {
			return nil, fmt.Errorf("invalid tournament ID format")
		}
		pipeline = append(pipeline, bson.M{"$match": bson.M{"tournament_id": objID}})
	}

	// Add stages to lookup team details
	lookupStages := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "teams",
				"localField":   "team_a_id",
				"foreignField": "_id",
				"as":           "team_a_details",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "teams",
				"localField":   "team_b_id",
				"foreignField": "_id",
				"as":           "team_b_details",
			},
		},
		{
			"$project": bson.M{
				"_id":                 1,
				"tournament_id":       1,
				"team_a_id":           1,
				"team_b_id":           1,
				"match_date":          1,
				"match_time":          1,
				"location":            1,
				"round":               1,
				"result_team_a_score": 1,
				"result_team_b_score": 1,
				"winner_team_id":      1,
				"status":              1,
				"created_at":          1,
				"updated_at":          1,
				"team_a": bson.M{
					"$arrayElemAt": []interface{}{
						bson.M{
							"$map": bson.M{
								"input": "$team_a_details", "as": "team",
								"in": bson.M{"_id": "$$team._id", "team_name": "$$team.team_name", "logo_url": "$$team.logo_url"},
							},
						}, 0,
					},
				},
				"team_b": bson.M{
					"$arrayElemAt": []interface{}{
						bson.M{
							"$map": bson.M{
								"input": "$team_b_details", "as": "team",
								"in": bson.M{"_id": "$$team._id", "team_name": "$$team.team_name", "logo_url": "$$team.logo_url"},
							},
						}, 0,
					},
				},
			},
		},
	}

	pipeline = append(pipeline, lookupStages...)

	cursor, err := config.MatchesCollection.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println("GetAllMatches (Aggregate):", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var matches []model.MatchWithDetails
	if err := cursor.All(ctx, &matches); err != nil {
		fmt.Println("GetAllMatches (Decode):", err)
		return nil, err
	}

	return matches, nil
}

// GetMatchByID retrieves match by ID
func GetMatchByID(ctx context.Context, id string) (*model.Match, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid match ID format")
	}

	var match model.Match
	filter := bson.M{"_id": objID}
	err = config.MatchesCollection.FindOne(ctx, filter).Decode(&match)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("terjadi kesalahan dalam mengambil data: %v", err)
	}
	return &match, nil
}

// GetMatchWithDetailsByID retrieves match by ID with team details populated
func GetMatchWithDetailsByID(ctx context.Context, id string) (*model.MatchWithDetails, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid match ID format")
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": objID},
		},
		{
			"$lookup": bson.M{
				"from":         "teams",
				"localField":   "team_a_id",
				"foreignField": "_id",
				"as":           "team_a_details",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "teams",
				"localField":   "team_b_id",
				"foreignField": "_id",
				"as":           "team_b_details",
			},
		},
		{
			"$project": bson.M{
				"_id":                 1,
				"tournament_id":       1,
				"team_a_id":           1,
				"team_b_id":           1,
				"match_date":          1,
				"match_time":          1,
				"location":            1,
				"round":               1,
				"result_team_a_score": 1,
				"result_team_b_score": 1,
				"winner_team_id":      1,
				"status":              1,
				"created_at":          1,
				"updated_at":          1,
				"team_a": bson.M{
					"$arrayElemAt": []interface{}{
						bson.M{
							"$map": bson.M{
								"input": "$team_a_details",
								"as":    "team",
								"in": bson.M{
									"_id":       "$$team._id",
									"team_name": "$$team.team_name",
									"logo_url":  "$$team.logo_url",
								},
							},
						}, 0,
					},
				},
				"team_b": bson.M{
					"$arrayElemAt": []interface{}{
						bson.M{
							"$map": bson.M{
								"input": "$team_b_details",
								"as":    "team",
								"in": bson.M{
									"_id":       "$$team._id",
									"team_name": "$$team.team_name",
									"logo_url":  "$$team.logo_url",
								},
							},
						}, 0,
					},
				},
			},
		},
	}

	cursor, err := config.MatchesCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("terjadi kesalahan dalam mengambil data: %v", err)
	}
	defer cursor.Close(ctx)

	var matches []model.MatchWithDetails
	if err := cursor.All(ctx, &matches); err != nil {
		return nil, fmt.Errorf("terjadi kesalahan dalam decode data: %v", err)
	}

	if len(matches) == 0 {
		return nil, nil
	}

	return &matches[0], nil
}

// UpdateMatch updates match data
func UpdateMatch(ctx context.Context, id string, update bson.M) (updatedID string, err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", fmt.Errorf("invalid match ID format")
	}

	// Basic validation for team IDs if they exist in the update map
	if teamAID, ok := update["team_a_id"]; ok {
		if _, ok := teamAID.(primitive.ObjectID); !ok {
			return "", fmt.Errorf("team_a_id must be a valid ObjectID")
		}
	}
	if teamBID, ok := update["team_b_id"]; ok {
		if _, ok := teamBID.(primitive.ObjectID); !ok {
			return "", fmt.Errorf("team_b_id must be a valid ObjectID")
		}
	}

	// Validate teams are different if both are being updated
	teamAID, teamAOK := update["team_a_id"]
	teamBID, teamBOK := update["team_b_id"]
	if teamAOK && teamBOK && teamAID == teamBID {
		return "", fmt.Errorf("Team A dan Team B harus berbeda")
	}

	// Note: More complex validations like checking if team IDs exist in the teams collection,
	// or validating winner_team_id against the participants might be needed here if not handled at the handler level.
	// For now, we assume handler-level validation is sufficient for partial updates.

	update["updated_at"] = time.Now()

	filter := bson.M{"_id": objID}
	updateData := bson.M{"$set": update}

	result, err := config.MatchesCollection.UpdateOne(ctx, filter, updateData)
	if err != nil {
		fmt.Printf("UpdateMatch: %v\n", err)
		return "", err
	}
	if result.ModifiedCount == 0 {
		// This could also mean the data sent was the same as the existing data.
		// To differentiate, a find and compare would be needed, but for now, this is acceptable.
		return "", fmt.Errorf("tidak ada data yang diupdate untuk Match ID %s, atau data yang dikirim sama", id)
	}
	return id, nil
}

// DeleteMatch deletes match by ID
func DeleteMatch(ctx context.Context, id string) (deletedID string, err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", fmt.Errorf("invalid match ID format")
	}

	filter := bson.M{"_id": objID}
	result, err := config.MatchesCollection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Printf("DeleteMatch: %v\n", err)
		return "", err
	}
	if result.DeletedCount == 0 {
		return "", fmt.Errorf("tidak ada data yang dihapus untuk Match ID %s", id)
	}
	return id, nil
}
