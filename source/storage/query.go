package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var timeout = 5 * time.Second

func InsertMany(collectionName string, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	collection := GetCollection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return collection.InsertMany(ctx, documents, opts...)
}

func FindOne(collectionName string, document interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	collection := GetCollection(collectionName)

	doc, err := convertToBsonM(document)
	if err != nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return collection.FindOne(ctx, doc, opts...)
}
