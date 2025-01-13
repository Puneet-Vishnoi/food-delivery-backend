package db

import (
	"context"
	"os"
	"time"

	constant "github.com/MarNawar/food-delivery-backend/user-service/internal/domain/entity/models/constants"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Db struct {
	MongoClient *mongo.Client
}

func ConnectDB() Db {
	// Create an empty background context to keep alive till the application is running.
	ctx := context.Background()

	var err error
	var mongoClient *mongo.Client

	// Try connecting to the database for the defined number of attempts.
	for i := 0; i < constant.MAX_DB_ATTEMPTS; i++ {

		// Attempt a client establishment.
		mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv(constant.MDBUri)))
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		// If the execution reaches here, the connection is successful.
		// Otherwise, the retires would result in loop reiterating till max attempts are crossed and loop is exited.
		break
	}

	if err == nil {

		// Return MongoDB client and associated data.

		err = mongoClient.Ping(context.Background(), nil)
		// In case of an error, return an invalid empty Server object to trigger a graceful shutdown.
		if err != nil {
			return Db{}
		} else {
		}
		return Db{MongoClient: mongoClient}

	} else {
		// If connection could not be established, return nil.
		return Db{}
	}
}

func (db *Db) Stop() {
	ctx := context.Background()
	err := db.MongoClient.Disconnect(ctx)
	if err != nil {
	}
}
