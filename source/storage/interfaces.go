package storage

const VideoMetadataC = "video_metadata"

type VideoMetadataInterface interface {
	BulkInsertMetadata(videoMetadatas []*VideoMetadata) error
	MetadataExists(id string) (bool, error)
	FindLastInsertedIndex() (int, error)
	FindOneMetadata(title string, description string) (*VideoMetadata, error)
}
