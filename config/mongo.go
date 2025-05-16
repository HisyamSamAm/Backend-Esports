package config

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoString string = os.Getenv("MONGOS")
var DBName = "mpls14"

var DB *mongo.Database

func MongoConnect(dbname string) (db *mongo.Database) {
    clientOpts := options.Client().ApplyURI(MongoString)

    client, err := mongo.Connect(context.TODO(), clientOpts)
    if err != nil {
        fmt.Println("MongoConnect: failed to connect:", err)
        return nil
    }

    // Ping untuk pastikan koneksi sukses
    if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
        fmt.Println("MongoConnect: ping failed:", err)
        return nil
    }

    fmt.Println("MongoConnect: connected to MongoDB Atlas")
    return client.Database(dbname)
}