package main

import (
	"context"
	"embeck/config"
	"encoding/json"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	// Memuat .env dari direktori root proyek
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Peringatan: Tidak dapat memuat file .env. Pastikan file tersebut ada di direktori root.")
	}
}

// prettyPrint mencetak data dalam format JSON yang mudah dibaca.
func prettyPrint(label string, data interface{}) {
	fmt.Printf("\n--- %s ---\n", label)
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Error saat marshalling data: %v", err)
		fmt.Println(data)
		return
	}
	fmt.Println(string(b))
}

func main() {
	fmt.Println("Memulai skrip diagnostik tim...")

	// --- KONFIGURASI ---
	// ID untuk tim dan kapten yang bermasalah ("dawdwa" dan "Tes 1")
	teamIDStr := "687e0937eee108e4f1995833"
	captainIDStr := "687e08a6eee108e4f1995832"

	// Hubungkan ke MongoDB
	db := config.MongoConnect(config.DBName)
	if db == nil {
		log.Fatal("Gagal terhubung ke database. Keluar.")
	}
	fmt.Println("✅ Berhasil terhubung ke database.")

	ctx := context.Background()
	teamObjID, err := primitive.ObjectIDFromHex(teamIDStr)
	if err != nil {
		log.Fatalf("Format ID Tim tidak valid: %v", err)
	}
	captainObjID, err := primitive.ObjectIDFromHex(captainIDStr)
	if err != nil {
		log.Fatalf("Format ID Kapten tidak valid: %v", err)
	}

	// Langkah 1: Temukan dokumen tim berdasarkan ID-nya
	var teamResult bson.M
	err = config.TeamsCollection.FindOne(ctx, bson.M{"_id": teamObjID}).Decode(&teamResult)
	if err != nil {
		log.Fatalf("Langkah 1 GAGAL: Tidak dapat menemukan tim dengan ID %s. Error: %v", teamIDStr, err)
	}
	prettyPrint("Langkah 1 BERHASIL: Dokumen Tim Ditemukan", teamResult)

	// Langkah 2: Temukan dokumen pemain (kapten) berdasarkan ID-nya
	var playerResult bson.M
	err = config.PlayersCollection.FindOne(ctx, bson.M{"_id": captainObjID}).Decode(&playerResult)
	if err != nil {
		log.Fatalf("Langkah 2 GAGAL: Tidak dapat menemukan pemain dengan ID %s. Error: %v", captainIDStr, err)
	}
	prettyPrint("Langkah 2 BERHASIL: Dokumen Pemain (Kapten) Ditemukan", playerResult)

	// Langkah 3: Jalankan pipeline agregasi untuk tim spesifik ini
	pipeline := []bson.M{
		{"$match": bson.M{"_id": teamObjID}},
		{
			"$lookup": bson.M{
				"from":         "players",
				"localField":   "captain_id",
				"foreignField": "_id",
				"as":           "captain_details_lookup",
			},
		},
	}

	cursor, err := config.TeamsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatalf("Langkah 3 GAGAL: Kueri agregasi gagal. Error: %v", err)
	}

	var aggregationResult []bson.M
	if err = cursor.All(ctx, &aggregationResult); err != nil {
		log.Fatalf("Langkah 3 GAGAL: Tidak dapat mendekode hasil agregasi. Error: %v", err)
	}

	if len(aggregationResult) == 0 {
		log.Fatalf("Langkah 3 GAGAL: Agregasi tidak mengembalikan dokumen apa pun.")
	}

	prettyPrint("Langkah 3 BERHASIL: Hasil Agregasi", aggregationResult[0])

	// Pemeriksaan akhir pada hasil lookup
	captainDetails := aggregationResult[0]["captain_details_lookup"]
	if details, ok := captainDetails.(primitive.A); ok && len(details) > 0 {
		fmt.Println("\n✅ DIAGNOSIS: SUKSES! Operasi $lookup bekerja dengan benar untuk tim ini.")
	} else {
		fmt.Println("\n❌ DIAGNOSIS: GAGAL! Operasi $lookup mengembalikan array kosong. Ini adalah inti masalahnya.")
	}
}
