package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoString string = os.Getenv("MONGOS")
var DBName = "mpls14"

var DB *mongo.Database

func ConnectDB() {
	mongoString := os.Getenv("MONGOS")
	if mongoString == "" {
		log.Fatal("Error connecting to MongoDB, dimana envnya!")
	}

	clientOpts := options.Client().ApplyURI(mongoString).
	SetServerSelectionTimeout(5 * time.Second)

	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}

	DB = client.Database(DBName)
	fmt.Println("Asikk connect nichh!")

}
