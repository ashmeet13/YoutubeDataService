package storage

import "time"

const VideoMetadataC = "video_metadata"

type VideoMetadata struct {
	VideoID              string    `bson:"video_id"`
	Title                string    `bson:"title"`
	Description          string    `bson:"description"`
	DefaultThumbnailURL  string    `bson:"default_thumbnail_url"`
	HighThumbnailURL     string    `bson:"high_thumbnail_url"`
	MaxresThumbnailURL   string    `bson:"maxres_thumbnail_url"`
	MediumThumbnailURL   string    `bson:"medium_thumbnail_url"`
	StandardThumbnailURL string    `bson:"standard_thumbnail_url"`
	PublishedAt          time.Time `bson:"published_at"`
}

const UserC = "users"

type User struct {
	UserID    string    `bson:"user_id"`
	PageSize  int       `bson:"page_size"`
	Timestamp time.Time `bson:"timestamp"`
}
