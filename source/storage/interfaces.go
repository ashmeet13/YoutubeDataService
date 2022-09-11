package storage

import "time"

type VideoMetadataInterface interface {
	BulkInsertMetadata(videoMetadatas []*VideoMetadata) error
	FindOneMetadataWithVideoID(id string) (*VideoMetadata, error)
	FindLastInsertedMetadata() (*VideoMetadata, error)
	FindOneMetadata(title string, description string) (*VideoMetadata, error)
	UpdateOneMetadata(id string, videoMetadata *VideoMetadata) error
	FetchPagedMetadata(timestamp time.Time, offset, limit int64) ([]*VideoMetadata, error)
}

type UserInterface interface {
	CreateUser(user *User) error
	ReadUser(userID string) (*User, error)
	UpdateUser(id string, user *User) error
}
