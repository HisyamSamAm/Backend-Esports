package main

import (
	"context"
	"embeck/config"
	"embeck/model"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Peringatan: Tidak dapat memuat file .env. Pastikan file tersebut ada di direktori root.")
	}
}

func main() {
	fmt.Println("Memulai proses seeding data dummy...")

	db := config.MongoConnect(config.DBName)
	if db == nil {
		log.Fatal("Gagal terhubung ke database. Keluar.")
	}
	fmt.Println("✅ Berhasil terhubung ke database.")

	ctx := context.Background()

	// 1. Hapus data lama untuk memastikan kebersihan data
	fmt.Println("\n--- Menghapus Data Lama ---")
	collectionsToClear := []*mongo.Collection{
		config.MatchesCollection,
		config.TournamentsCollection,
		config.TeamsCollection,
		config.PlayersCollection,
	}
	for _, coll := range collectionsToClear {
		_, err := coll.DeleteMany(ctx, primitive.M{})
		if err != nil {
			log.Fatalf("Gagal membersihkan koleksi %s: %v", coll.Name(), err)
		}
		fmt.Printf("Koleksi '%s' berhasil dibersihkan.\n", coll.Name())
	}

	// 2. Buat Data Pemain (Players)
	fmt.Println("\n--- Membuat Data Pemain ---")
	players := generatePlayers()
	playerIDs, err := insertPlayers(ctx, players)
	if err != nil {
		log.Fatalf("Gagal membuat data pemain: %v", err)
	}
	fmt.Printf("✅ Berhasil membuat %d pemain.\n", len(playerIDs))

	// 3. Buat Data Tim (Teams)
	fmt.Println("\n--- Membuat Data Tim ---")
	teams := generateTeams(playerIDs)
	teamIDs, err := insertTeams(ctx, teams)
	if err != nil {
		log.Fatalf("Gagal membuat data tim: %v", err)
	}
	fmt.Printf("✅ Berhasil membuat %d tim.\n", len(teamIDs))

	// 4. Buat Data Turnamen (Tournament)
	fmt.Println("\n--- Membuat Data Turnamen ---")
	tournament := generateTournament(teamIDs)
	tournamentID, err := insertTournament(ctx, tournament)
	if err != nil {
		log.Fatalf("Gagal membuat data turnamen: %v", err)
	}
	fmt.Printf("✅ Berhasil membuat turnamen '%s'.\n", tournament.Name)

	// 5. Buat Data Pertandingan (Matches)
	fmt.Println("\n--- Membuat Data Pertandingan ---")
	matches := generateMatches(tournamentID, teamIDs)
	matchCount, err := insertMatches(ctx, matches)
	if err != nil {
		log.Fatalf("Gagal membuat data pertandingan: %v", err)
	}
	fmt.Printf("✅ Berhasil membuat %d pertandingan.\n", matchCount)

	fmt.Println("\n🎉 Proses seeding data dummy selesai!")
}

// --- Fungsi Helper untuk Generate dan Insert Data ---

func generatePlayers() []model.Player {
	return []model.Player{
		{ID: primitive.NewObjectID(), Name: "Muhammad Ikhsan", MLNickname: "Lemon", MLID: "100001", Status: "Active"},
		{ID: primitive.NewObjectID(), Name: "Calvin Winata", MLNickname: "Vyn", MLID: "100002", Status: "Active"},
		{ID: primitive.NewObjectID(), Name: "Albert Neilsen Iskandar", MLNickname: "Alberttt", MLID: "100003", Status: "Active"},
		{ID: primitive.NewObjectID(), Name: "Maxhill Leonardo", MLNickname: "Antimage", MLID: "100004", Status: "Active"},
		{ID: primitive.NewObjectID(), Name: "Ferdyansyah Kamaruddin", MLNickname: "Ferxiic", MLID: "100005", Status: "Active"},
		{ID: primitive.NewObjectID(), Name: "Gilang", MLNickname: "Sanz", MLID: "100006", Status: "Active"},
		{ID: primitive.NewObjectID(), Name: "Nicky Fernando", MLNickname: "Kiboy", MLID: "100007", Status: "Active"},
		{ID: primitive.NewObjectID(), Name: "Calvin", MLNickname: "CW", MLID: "100008", Status: "Active"},
		{ID: primitive.NewObjectID(), Name: "Kairi Rayosdelsol", MLNickname: "Kairi", MLID: "100009", Status: "Active"},
		{ID: primitive.NewObjectID(), Name: "Buts", MLNickname: "Butsss", MLID: "100010", Status: "Active"},
		{ID: primitive.NewObjectID(), Name: "Ihsan Besari Kusudana", MLNickname: "Luminaire", MLID: "100011", Status: "Inactive"},
		{ID: primitive.NewObjectID(), Name: "Eko Julianto", MLNickname: "Oura", MLID: "100012", Status: "Inactive"},
	}
}

func insertPlayers(ctx context.Context, players []model.Player) (map[string]primitive.ObjectID, error) {
	playerDocs := make([]interface{}, len(players))
	playerIDs := make(map[string]primitive.ObjectID)
	for i, p := range players {
		p.CreatedAt = time.Now()
		p.UpdatedAt = time.Now()
		playerDocs[i] = p
		playerIDs[p.MLNickname] = p.ID
	}
	_, err := config.PlayersCollection.InsertMany(ctx, playerDocs)
	return playerIDs, err
}

func generateTeams(playerIDs map[string]primitive.ObjectID) []model.Team {
	return []model.Team{
		{
			ID:        primitive.NewObjectID(),
			TeamName:  "RRQ",
			CaptainID: playerIDs["Lemon"],
			Members:   []primitive.ObjectID{playerIDs["Lemon"], playerIDs["Vyn"], playerIDs["Alberttt"], playerIDs["Antimage"]},
		},
		{
			ID:        primitive.NewObjectID(),
			TeamName:  "ONIC Esports",
			CaptainID: playerIDs["Sanz"],
			Members:   []primitive.ObjectID{playerIDs["Sanz"], playerIDs["Kiboy"], playerIDs["CW"], playerIDs["Kairi"], playerIDs["Butsss"]},
		},
		{
			ID:        primitive.NewObjectID(),
			TeamName:  "EVOS Legends",
			CaptainID: playerIDs["Luminaire"],
			Members:   []primitive.ObjectID{playerIDs["Luminaire"], playerIDs["Ferxiic"], playerIDs["Oura"]},
		},
	}
}

func insertTeams(ctx context.Context, teams []model.Team) (map[string]primitive.ObjectID, error) {
	teamDocs := make([]interface{}, len(teams))
	teamIDs := make(map[string]primitive.ObjectID)
	for i, t := range teams {
		t.CreatedAt = time.Now()
		t.UpdatedAt = time.Now()
		teamDocs[i] = t
		teamIDs[t.TeamName] = t.ID
	}
	_, err := config.TeamsCollection.InsertMany(ctx, teamDocs)
	return teamIDs, err
}

func generateTournament(teamIDs map[string]primitive.ObjectID) model.Tournament {
	return model.Tournament{
		ID:                 primitive.NewObjectID(),
		Name:               "MPL Indonesia Season 15",
		Description:        "Turnamen Mobile Legends paling bergengsi di Indonesia.",
		StartDate:          time.Now().AddDate(0, 1, 0), // Mulai bulan depan
		EndDate:            time.Now().AddDate(0, 3, 0), // Berakhir 3 bulan dari sekarang
		PrizePool:          "$300,000 USD",
		Status:             "upcoming",
		TeamsParticipating: []primitive.ObjectID{teamIDs["RRQ"], teamIDs["ONIC Esports"], teamIDs["EVOS Legends"]},
		CreatedBy:          primitive.NewObjectID(), // Dummy Admin ID
	}
}

func insertTournament(ctx context.Context, t model.Tournament) (primitive.ObjectID, error) {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	_, err := config.TournamentsCollection.InsertOne(ctx, t)
	return t.ID, err
}

func generateMatches(tournamentID primitive.ObjectID, teamIDs map[string]primitive.ObjectID) []model.Match {
	return []model.Match{
		{ // Match 1: RRQ vs ONIC
			TournamentID: tournamentID,
			TeamAID:      teamIDs["RRQ"],
			TeamBID:      teamIDs["ONIC Esports"],
			MatchDate:    time.Now().AddDate(0, 1, 7), // Minggu pertama turnamen
			MatchTime:    "18:00 WIB",
			Round:        "Regular Season - Week 1",
			Status:       "scheduled",
		},
		{ // Match 2: EVOS vs RRQ
			TournamentID: tournamentID,
			TeamAID:      teamIDs["EVOS Legends"],
			TeamBID:      teamIDs["RRQ"],
			MatchDate:    time.Now().AddDate(0, 1, 8),
			MatchTime:    "20:00 WIB",
			Round:        "Regular Season - Week 1",
			Status:       "scheduled",
		},
		{ // Match 3: ONIC vs EVOS
			TournamentID: tournamentID,
			TeamAID:      teamIDs["ONIC Esports"],
			TeamBID:      teamIDs["EVOS Legends"],
			MatchDate:    time.Now().AddDate(0, 1, 14), // Minggu kedua
			MatchTime:    "18:00 WIB",
			Round:        "Regular Season - Week 2",
			Status:       "scheduled",
		},
	}
}

func insertMatches(ctx context.Context, matches []model.Match) (int, error) {
	matchDocs := make([]interface{}, len(matches))
	for i, m := range matches {
		m.ID = primitive.NewObjectID()
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()
		matchDocs[i] = m
	}
	_, err := config.MatchesCollection.InsertMany(ctx, matchDocs)
	return len(matches), err
}
