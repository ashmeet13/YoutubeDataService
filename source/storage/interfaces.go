package storage

const VideoMetadataC = "video_metadata"

type VideoMetadataInterface interface {
	BulkInsertMetadata(videoMetadatas []*VideoMetadata) error
	FindOneMetadataWithVideoID(id string) (*VideoMetadata, error)
	FindLastInsertedMetadata() (*VideoMetadata, error)
	FindOneMetadata(title string, description string) (*VideoMetadata, error)
	UpdateOneMetadata(id string, videoMetadata *VideoMetadata) error
	FetchPage(start, end int) ([]*VideoMetadata, error)
}
