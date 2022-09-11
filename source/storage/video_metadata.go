package storage

import (
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

func (m *VideoMetadataImpl) MetadataExists(id string) (bool, error) {
	query := bson.M{
		"video_id": bson.M{"$eq": id},
	}

	result := FindOne(m.collection, query)

	var decodedResult VideoMetadata
	err := result.Decode(&decodedResult)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, err
}

func (m *VideoMetadataImpl) FindLastInsertedIndex() (int, error) {
	result := FindOne(m.collection, bson.M{}, &options.FindOneOptions{Sort: bson.M{"document_index": -1}})

	var decodedResult VideoMetadata
	err := result.Decode(&decodedResult)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return -1, nil
		}
		return 0, err
	}

	return decodedResult.DocumentIndex, err
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
