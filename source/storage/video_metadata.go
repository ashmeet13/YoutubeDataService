package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewVideoMetadataImpl() *VideoMetadataImpl {
	return &VideoMetadataImpl{
		collection: VideoMetadataC,
	}
}

type VideoMetadataImpl struct {
	collection string
}

func (m *VideoMetadataImpl) BulkInsertMetadata(videoMetadatas []*VideoMetadata) error {
	insertDocs := bson.A{}

	for _, metadata := range videoMetadatas {
		doc, err := convertToBsonM(metadata)
		if err != nil {
			return err
		}

		insertDocs = append(insertDocs, doc)
	}

	_, err := InsertMany(m.collection, insertDocs)
	if err != nil {
		return err
	}

	return nil
}

func (m *VideoMetadataImpl) FindOneMetadataWithVideoID(id string) (*VideoMetadata, error) {
	query := bson.M{
		"video_id": bson.M{"$eq": id},
	}

	result := FindOne(m.collection, query)

	var decodedResult VideoMetadata
	err := result.Decode(&decodedResult)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &decodedResult, err
}

func (m *VideoMetadataImpl) FindLastInsertedMetadata() (*VideoMetadata, error) {
	result := FindOne(m.collection, bson.M{}, &options.FindOneOptions{Sort: bson.M{"document_index": -1}})

	var decodedResult VideoMetadata
	err := result.Decode(&decodedResult)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &decodedResult, err
}

func (m *VideoMetadataImpl) FindOneMetadata(title string, description string) (*VideoMetadata, error) {
	query := bson.M{}

	if title != "" {
		query["title"] = bson.M{"$eq": title}
	}

	if description != "" {
		query["description"] = bson.M{"$eq": description}
	}

	result := FindOne(m.collection, query)

	var decodedResult VideoMetadata
	err := result.Decode(&decodedResult)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &decodedResult, nil
}

func (m *VideoMetadataImpl) UpdateOneMetadata(id string, videoMetadata *VideoMetadata) error {
	filters := bson.M{
		"video_id": bson.M{"$eq": id},
	}

	modifier := bson.M{
		"$set": videoMetadata,
	}

	_, err := UpdateOne(m.collection, filters, modifier)
	if err != nil {
		return err
	}

	return nil
}

func (m *VideoMetadataImpl) FetchPage(start, end int) ([]*VideoMetadata, error) {
	query := bson.M{
		"document_index": bson.M{
			"$gt":  start,
			"$lte": end,
		},
	}

	cur, err := Find(m.collection, query)

	err = cur.Err()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	defer cur.Close(ctx)

	var metadata []*VideoMetadata
	for cur.Next(ctx) {
		var videoMetadata VideoMetadata
		err := cur.Decode(&videoMetadata)
		if err != nil {
			return nil, err
		}
		metadata = append(metadata, &videoMetadata)
	}

	return metadata, nil
}
