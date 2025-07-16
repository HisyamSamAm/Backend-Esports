package main

import (
	"context"
	"embeck/config"
	"embeck/model"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	// Load .env file if exists
	if _, err := os.Stat("../../.env"); err == nil {
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Println("Warning: Could not load .env file")
		}
	}
}

func main() {
	// Connect to database
	db := config.MongoConnect(config.DBName)
	if db == nil {
		log.Fatal("Failed to connect to database")
	}

	fmt.Println("🎮 Creating 60 unique players (5 players for each of 12 teams)...")

	ctx := context.Background()

	// Create 60 unique players - organized by teams
	players := []model.Player{
		// Tournament 1 Teams (6 teams x 5 players = 30 players)

		// Team 1: RRQ Hoshi
		{ID: primitive.NewObjectID(), Name: "Ahmad Febriyanto", MLNickname: "Lemon", MLID: "12345678901", Status: "active", AvatarURL: "rrq_lemon.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Rivaldi Dwi Saputra", MLNickname: "R7", MLID: "12345678902", Status: "active", AvatarURL: "rrq_r7.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Calvin Winata", MLNickname: "Vyn", MLID: "12345678903", Status: "active", AvatarURL: "rrq_vyn.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Muhammad Gilang", MLNickname: "Skylar", MLID: "12345678904", Status: "active", AvatarURL: "rrq_skylar.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Adit Rahman", MLNickname: "Xinnn", MLID: "12345678905", Status: "active", AvatarURL: "rrq_xinnn.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},

		// Team 2: EVOS Legends
		{ID: primitive.NewObjectID(), Name: "Wawan Suherman", MLNickname: "Wann", MLID: "12345678906", Status: "active", AvatarURL: "evos_wann.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Antimage Dewa", MLNickname: "Antimage", MLID: "12345678907", Status: "active", AvatarURL: "evos_antimage.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Ferxiic Pratama", MLNickname: "Ferxiic", MLID: "12345678908", Status: "active", AvatarURL: "evos_ferxiic.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Luminaire Gaming", MLNickname: "Luminaire", MLID: "12345678909", Status: "active", AvatarURL: "evos_luminaire.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Clover Gaming", MLNickname: "Clover", MLID: "12345678910", Status: "active", AvatarURL: "evos_clover.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},

		// Team 3: Bigetron Alpha
		{ID: primitive.NewObjectID(), Name: "Branz Valentino", MLNickname: "Branz", MLID: "12345678911", Status: "active", AvatarURL: "btr_branz.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Zxuan Gaming", MLNickname: "Zxuan", MLID: "12345678912", Status: "active", AvatarURL: "btr_zxuan.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Dlar Kurniawan", MLNickname: "Dlar", MLID: "12345678913", Status: "active", AvatarURL: "btr_dlar.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Mikoto Gaming", MLNickname: "Mikoto", MLID: "12345678914", Status: "active", AvatarURL: "btr_mikoto.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Dominus Gaming", MLNickname: "Dominus", MLID: "12345678915", Status: "active", AvatarURL: "btr_dominus.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},

		// Team 4: ONIC Esports
		{ID: primitive.NewObjectID(), Name: "Sanz Jangkar", MLNickname: "Sanz", MLID: "12345678916", Status: "active", AvatarURL: "onic_sanz.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "CW Hero", MLNickname: "CW", MLID: "12345678917", Status: "active", AvatarURL: "onic_cw.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Kairi Legend", MLNickname: "Kairi", MLID: "12345678918", Status: "active", AvatarURL: "onic_kairi.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Butsss Gaming", MLNickname: "Butsss", MLID: "12345678919", Status: "active", AvatarURL: "onic_butsss.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Kiboy Gaming", MLNickname: "Kiboy", MLID: "12345678920", Status: "active", AvatarURL: "onic_kiboy.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},

		// Team 5: Geek Fam
		{ID: primitive.NewObjectID(), Name: "Jess No Limit", MLNickname: "JessNoLimit", MLID: "12345678921", Status: "active", AvatarURL: "gf_jess.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Skylar Gaming", MLNickname: "Skylar", MLID: "12345678922", Status: "active", AvatarURL: "gf_skylar.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Udil Mabar", MLNickname: "Udil", MLID: "12345678923", Status: "active", AvatarURL: "gf_udil.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Celiboy Gaming", MLNickname: "Celiboy", MLID: "12345678924", Status: "active", AvatarURL: "gf_celiboy.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Nexus Gaming", MLNickname: "Nexus", MLID: "12345678925", Status: "active", AvatarURL: "gf_nexus.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},

		// Team 6: Aura Fire
		{ID: primitive.NewObjectID(), Name: "Oura Gaming", MLNickname: "Oura", MLID: "12345678926", Status: "active", AvatarURL: "af_oura.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Celiboy Gaming", MLNickname: "Celiboy", MLID: "12345678927", Status: "active", AvatarURL: "af_celiboy.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Clay Gaming", MLNickname: "Clay", MLID: "12345678928", Status: "active", AvatarURL: "af_clay.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Rekt Gaming", MLNickname: "Rekt", MLID: "12345678929", Status: "active", AvatarURL: "af_rekt.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Wannnn Gaming", MLNickname: "Wannnn", MLID: "12345678930", Status: "active", AvatarURL: "af_wannnn.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},

		// Tournament 2 Teams (6 teams x 5 players = 30 players)

		// Team 7: Rebellion Zion
		{ID: primitive.NewObjectID(), Name: "Tuanmuda Gaming", MLNickname: "Tuanmuda", MLID: "12345678931", Status: "active", AvatarURL: "rbz_tuanmuda.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Vyncent Gaming", MLNickname: "Vyncent", MLID: "12345678932", Status: "active", AvatarURL: "rbz_vyncent.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Clayyy Gaming", MLNickname: "Clayyy", MLID: "12345678933", Status: "active", AvatarURL: "rbz_clayyy.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Aboy Gaming", MLNickname: "Aboy", MLID: "12345678934", Status: "active", AvatarURL: "rbz_aboy.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Drakath Gaming", MLNickname: "Drakath", MLID: "12345678935", Status: "active", AvatarURL: "rbz_drakath.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},

		// Team 8: GPX Basreng
		{ID: primitive.NewObjectID(), Name: "Potato Gaming", MLNickname: "Potato", MLID: "12345678936", Status: "active", AvatarURL: "gpx_potato.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Maxxx Gaming", MLNickname: "Maxxx", MLID: "12345678937", Status: "active", AvatarURL: "gpx_maxxx.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Fierce Gaming", MLNickname: "Fierce", MLID: "12345678938", Status: "active", AvatarURL: "gpx_fierce.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Xinnn Gaming", MLNickname: "Xinnn", MLID: "12345678939", Status: "active", AvatarURL: "gpx_xinnn.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Clover Gaming", MLNickname: "Clover", MLID: "12345678940", Status: "active", AvatarURL: "gpx_clover.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},

		// Team 9: Alter Ego
		{ID: primitive.NewObjectID(), Name: "Ahmad Gaming", MLNickname: "Ahmad", MLID: "12345678941", Status: "active", AvatarURL: "ae_ahmad.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Celiboy Gaming", MLNickname: "Celiboy", MLID: "12345678942", Status: "active", AvatarURL: "ae_celiboy.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Rekt Gaming", MLNickname: "Rekt", MLID: "12345678943", Status: "active", AvatarURL: "ae_rekt.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Wannnn Gaming", MLNickname: "Wannnn", MLID: "12345678944", Status: "active", AvatarURL: "ae_wannnn.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Skylar Gaming", MLNickname: "Skylar", MLID: "12345678945", Status: "active", AvatarURL: "ae_skylar.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},

		// Team 10: Todak
		{ID: primitive.NewObjectID(), Name: "Moon Gaming", MLNickname: "Moon", MLID: "12345678946", Status: "active", AvatarURL: "tk_moon.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Yeb Gaming", MLNickname: "Yeb", MLID: "12345678947", Status: "active", AvatarURL: "tk_yeb.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Ejhay Gaming", MLNickname: "Ejhay", MLID: "12345678948", Status: "active", AvatarURL: "tk_ejhay.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Ly4 Gaming", MLNickname: "Ly4", MLID: "12345678949", Status: "active", AvatarURL: "tk_ly4.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Hesa Gaming", MLNickname: "Hesa", MLID: "12345678950", Status: "active", AvatarURL: "tk_hesa.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},

		// Team 11: Team Flash
		{ID: primitive.NewObjectID(), Name: "Flash Gaming", MLNickname: "Flash", MLID: "12345678951", Status: "active", AvatarURL: "tf_flash.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Speed Gaming", MLNickname: "Speed", MLID: "12345678952", Status: "active", AvatarURL: "tf_speed.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Lightning Gaming", MLNickname: "Lightning", MLID: "12345678953", Status: "active", AvatarURL: "tf_lightning.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Thunder Gaming", MLNickname: "Thunder", MLID: "12345678954", Status: "active", AvatarURL: "tf_thunder.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Storm Gaming", MLNickname: "Storm", MLID: "12345678955", Status: "active", AvatarURL: "tf_storm.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},

		// Team 12: Omega Esports
		{ID: primitive.NewObjectID(), Name: "Omega Gaming", MLNickname: "Omega", MLID: "12345678956", Status: "active", AvatarURL: "oe_omega.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Alpha Gaming", MLNickname: "Alpha", MLID: "12345678957", Status: "active", AvatarURL: "oe_alpha.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Beta Gaming", MLNickname: "Beta", MLID: "12345678958", Status: "active", AvatarURL: "oe_beta.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Gamma Gaming", MLNickname: "Gamma", MLID: "12345678959", Status: "active", AvatarURL: "oe_gamma.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: primitive.NewObjectID(), Name: "Delta Gaming", MLNickname: "Delta", MLID: "12345678960", Status: "active", AvatarURL: "oe_delta.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	// Insert all players
	var playerInterfaces []interface{}
	for _, player := range players {
		playerInterfaces = append(playerInterfaces, player)
	}

	result, err := config.PlayersCollection.InsertMany(ctx, playerInterfaces)
	if err != nil {
		log.Fatal("Failed to insert players:", err)
	}

	fmt.Printf("✅ Successfully created %d unique players!\n", len(result.InsertedIDs))
	fmt.Println("\n📋 Players organized by teams:")
	fmt.Println("   🏆 Tournament 1 Teams:")
	fmt.Println("     • RRQ Hoshi: Lemon, R7, Vyn, Skylar, Xinnn")
	fmt.Println("     • EVOS Legends: Wann, Antimage, Ferxiic, Luminaire, Clover")
	fmt.Println("     • Bigetron Alpha: Branz, Zxuan, Dlar, Mikoto, Dominus")
	fmt.Println("     • ONIC Esports: Sanz, CW, Kairi, Butsss, Kiboy")
	fmt.Println("     • Geek Fam: JessNoLimit, Skylar, Udil, Celiboy, Nexus")
	fmt.Println("     • Aura Fire: Oura, Celiboy, Clay, Rekt, Wannnn")
	fmt.Println("   🏆 Tournament 2 Teams:")
	fmt.Println("     • Rebellion Zion: Tuanmuda, Vyncent, Clayyy, Aboy, Drakath")
	fmt.Println("     • GPX Basreng: Potato, Maxxx, Fierce, Xinnn, Clover")
	fmt.Println("     • Alter Ego: Ahmad, Celiboy, Rekt, Wannnn, Skylar")
	fmt.Println("     • Todak: Moon, Yeb, Ejhay, Ly4, Hesa")
	fmt.Println("     • Team Flash: Flash, Speed, Lightning, Thunder, Storm")
	fmt.Println("     • Omega Esports: Omega, Alpha, Beta, Gamma, Delta")
	fmt.Println("\n✨ Each player is unique and will be assigned to exactly ONE team!")
}
