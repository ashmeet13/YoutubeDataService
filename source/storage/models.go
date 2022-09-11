package storage

import "time"

type VideoMetadata struct {
	VideoID              string    `bson:"video_id"`
	DocumentIndex        int       `bson:"document_index"`
	Title                string    `bson:"title"`
	Description          string    `bson:"description"`
	DefaultThumbnailURL  string    `bson:"default_thumbnail_url"`
	HighThumbnailURL     string    `bson:"high_thumbnail_url"`
	MaxresThumbnailURL   string    `bson:"maxres_thumbnail_url"`
	MediumThumbnailURL   string    `bson:"medium_thumbnail_url"`
	StandardThumbnailURL string    `bson:"standard_thumbnail_url"`
	PublishedAt          time.Time `bson:"published_at"`
}
