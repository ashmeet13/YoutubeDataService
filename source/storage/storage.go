package storage

import (
	"context"
	"time"

	"github.com/ashmeet13/YoutubeDataService/source/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client

func GetCollection(collectionName string) *mongo.Collection {
	database := GetDatabase()

	collection := database.Collection(collectionName)
	if collection == nil {
		panic("Collection nil in GetCollection")
	}
	return collection
}

func GetDatabase() *mongo.Database {
	config := common.GetConfiguration()
	if mongoClient == nil {
		initaliseMongoClient(config.MongoBaseURL)
	}

	return mongoClient.Database(config.MongoDatabaseName)
}

func Init() {
	config := common.GetConfiguration()
	initaliseMongoClient(config.MongoBaseURL)
}

func initaliseMongoClient(mongoBaseURL string) {
	logger := common.GetLogger()

	var err error
	mongoClientOptions := options.Client().ApplyURI(mongoBaseURL)

	ctx := context.Background()

	mongoClient, err = mongo.Connect(ctx, mongoClientOptions)
	if err != nil {
		panic("Unable to connect to mongo")
	}

	for {
		err = mongoClient.Ping(ctx, readpref.Primary())
		if err != nil {
			logger.Info("Unable to ping mongo server, sleeping for 5 seconds")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

}
