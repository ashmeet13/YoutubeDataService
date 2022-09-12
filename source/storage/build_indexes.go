package storage

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

var onceIndex sync.Once

func BuildIndexes(ctx context.Context) {
	onceIndex.Do(func() {
		db := GetDatabase()

		db.Collection(UserC).Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bsonx.Doc{
				{Key: "user_id", Value: bsonx.Int32(1)},
			},
			Options: options.Index().SetUnique(true).SetBackground(true),
		})

		db.Collection(VideoMetadataC).Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bsonx.Doc{
				{Key: "video_id", Value: bsonx.Int32(1)},
			},
			Options: options.Index().SetUnique(true).SetBackground(true),
		})

		db.Collection(VideoMetadataC).Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bsonx.Doc{
				{Key: "published_at", Value: bsonx.Int32(-1)},
			},
			Options: options.Index().SetUnique(true).SetBackground(true),
		})
	})
}
