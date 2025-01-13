package providers

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type DBHelper struct {
	MongoClient *mongo.Client
}

func NewDbProvider(mongoDBClient *mongo.Client) *DBHelper {
	return &DBHelper{
		MongoClient: mongoDBClient,
	}
}
