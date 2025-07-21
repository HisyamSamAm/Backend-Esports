package repository

import (
	"context"
	"embeck/config"
	"embeck/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateTournament creates a new tournament
func CreateTournament(tournament *model.Tournament) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tournament.CreatedAt = time.Now()
	tournament.UpdatedAt = time.Now()

	result, err := config.TournamentsCollection.InsertOne(ctx, tournament)
	return result, err
}

// GetAllTournaments retrieves all tournaments (admin view) with populated team details
func GetAllTournaments() ([]model.TournamentWithDetails, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "teams",
				"localField":   "teams_participating",
				"foreignField": "_id",
				"as":           "team_details",
			},
		},
		{
			"$project": bson.M{
				"_id":                1,
				"name":               1,
				"description":        1,
				"start_date":         1,
				"end_date":           1,
				"prize_pool":         1,
				"rules_document_url": 1,
				"status":             1,
				"created_by":         1,
				"created_at":         1,
				"updated_at":         1,
				"teams_participating": bson.M{
					"$map": bson.M{
						"input": "$team_details",
						"as":    "team",
						"in": bson.M{
							"_id":       "$$team._id",
							"team_name": "$$team.team_name",
							"logo_url":  "$$team.logo_url",
						},
					},
				},
				"matches": []bson.M{}, // Return empty array for matches to match the struct
			},
		},
	}

	cursor, err := config.TournamentsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tournaments []model.TournamentWithDetails
	if err = cursor.All(ctx, &tournaments); err != nil {
		return nil, err
	}

	return tournaments, nil
}

// GetAllTournamentsPublic retrieves all tournaments for public access (without admin fields)
func GetAllTournamentsPublic() ([]model.TournamentPublic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Project only public fields
	pipeline := []bson.M{
		{
			"$project": bson.M{
				"_id":                1,
				"name":               1,
				"description":        1,
				"start_date":         1,
				"end_date":           1,
				"prize_pool":         1,
				"rules_document_url": 1,
				"status":             1,
			},
		},
	}

	cursor, err := config.TournamentsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tournaments []model.TournamentPublic
	if err = cursor.All(ctx, &tournaments); err != nil {
		return nil, err
	}

	return tournaments, nil
}

// GetTournamentByID retrieves a tournament by ID (admin view)
func GetTournamentByID(id string) (*model.Tournament, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var tournament model.Tournament
	err = config.TournamentsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&tournament)
	if err != nil {
		return nil, err
	}

	return &tournament, nil
}

// GetTournamentWithDetailsByID retrieves tournament with populated teams and matches for public view
func GetTournamentWithDetailsByID(id string) (*model.TournamentWithDetails, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Aggregation pipeline to populate teams and get related matches
	pipeline := []bson.M{
		// Match the tournament
		{
			"$match": bson.M{"_id": objectID},
		},
		// Lookup teams participating
		{
			"$lookup": bson.M{
				"from":         "teams",
				"localField":   "teams_participating",
				"foreignField": "_id",
				"as":           "team_details",
			},
		},
		// Lookup matches for this tournament
		{
			"$lookup": bson.M{
				"from":         "matches",
				"localField":   "_id",
				"foreignField": "tournament_id",
				"as":           "match_details",
			},
		},
		// Lookup team A details for matches
		{
			"$lookup": bson.M{
				"from": "teams",
				"let":  bson.M{"match_details": "$match_details"},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$in": []interface{}{
									"$_id",
									bson.M{
										"$map": bson.M{
											"input": "$$match_details",
											"as":    "match",
											"in":    "$$match.team_a_id",
										},
									},
								},
							},
						},
					},
				},
				"as": "team_a_details",
			},
		},
		// Lookup team B details for matches
		{
			"$lookup": bson.M{
				"from": "teams",
				"let":  bson.M{"match_details": "$match_details"},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$in": []interface{}{
									"$_id",
									bson.M{
										"$map": bson.M{
											"input": "$$match_details",
											"as":    "match",
											"in":    "$$match.team_b_id",
										},
									},
								},
							},
						},
					},
				},
				"as": "team_b_details",
			},
		},
		// Project final structure
		{
			"$project": bson.M{
				"_id":                1,
				"name":               1,
				"description":        1,
				"start_date":         1,
				"end_date":           1,
				"prize_pool":         1,
				"rules_document_url": 1,
				"status":             1,
				"teams_participating": bson.M{
					"$map": bson.M{
						"input": "$team_details",
						"as":    "team",
						"in": bson.M{
							"_id":       "$$team._id",
							"team_name": "$$team.team_name",
							"logo_url":  "$$team.logo_url",
						},
					},
				},
				"matches": bson.M{
					"$map": bson.M{
						"input": "$match_details",
						"as":    "match",
						"in": bson.M{
							"_id":                 "$$match._id",
							"match_date":          "$$match.match_date",
							"match_time":          "$$match.match_time",
							"location":            "$$match.location",
							"round":               "$$match.round",
							"result_team_a_score": "$$match.result_team_a_score",
							"result_team_b_score": "$$match.result_team_b_score",
							"winner_team_id":      "$$match.winner_team_id",
							"status":              "$$match.status",
							"team_a": bson.M{
								"$let": bson.M{
									"vars": bson.M{
										"teamA": bson.M{
											"$arrayElemAt": []interface{}{
												bson.M{
													"$filter": bson.M{
														"input": "$team_a_details",
														"cond": bson.M{
															"$eq": []interface{}{"$$this._id", "$$match.team_a_id"},
														},
													},
												}, 0,
											},
										},
									},
									"in": bson.M{
										"_id":       "$$teamA._id",
										"team_name": "$$teamA.team_name",
										"logo_url":  "$$teamA.logo_url",
									},
								},
							},
							"team_b": bson.M{
								"$let": bson.M{
									"vars": bson.M{
										"teamB": bson.M{
											"$arrayElemAt": []interface{}{
												bson.M{
													"$filter": bson.M{
														"input": "$team_b_details",
														"cond": bson.M{
															"$eq": []interface{}{"$$this._id", "$$match.team_b_id"},
														},
													},
												}, 0,
											},
										},
									},
									"in": bson.M{
										"_id":       "$$teamB._id",
										"team_name": "$$teamB.team_name",
										"logo_url":  "$$teamB.logo_url",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	cursor, err := config.TournamentsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []model.TournamentWithDetails
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return &results[0], nil
}

// UpdateTournament updates a tournament
func UpdateTournament(id string, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update["updated_at"] = time.Now()

	_, err = config.TournamentsCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": update},
	)

	return err
}

// DeleteTournament deletes a tournament
func DeleteTournament(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = config.TournamentsCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// ValidateTeamsExist checks if all team IDs exist in the teams collection
func ValidateTeamsExist(teamIDs []string) error {
	if len(teamIDs) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectIDs := make([]primitive.ObjectID, len(teamIDs))
	for i, id := range teamIDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		objectIDs[i] = objectID
	}

	count, err := config.TeamsCollection.CountDocuments(ctx, bson.M{
		"_id": bson.M{"$in": objectIDs},
	})
	if err != nil {
		return err
	}

	if int(count) != len(teamIDs) {
		return mongo.ErrNoDocuments
	}

	return nil
}
