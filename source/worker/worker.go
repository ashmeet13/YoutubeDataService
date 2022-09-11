package worker

import (
	"strings"
	"time"

	"github.com/ashmeet13/YoutubeDataService/source/common"
	"github.com/ashmeet13/YoutubeDataService/source/storage"
	youtube_handler "github.com/ashmeet13/YoutubeDataService/source/youtube"
	"google.golang.org/api/youtube/v3"
)

/*
Worker is an async background worker that polls the Youtube API for Video events.

How does it work?

We query the GET /search/list endpoint to fetch data with the following parameters enabled
1. Parts = snippet, this fetches us the metadata we require
2. Type = video, we only want data for videos
3. OrderBy = date, we get the data ordered with most recent event at the start of the response
4. PublishedAfter = <time.RFC3339 formated date> we will be fetching all the data that was created on and after this point of time.

Assumption - /search/list can return data for a video which has been updated

We iterate over the result and create structs to insert into our DB.
While doing so we make a DB read call to check if the video already exists.

If it does and the published at of the new result is after the database doc, we update the doc for the video id.
Else, we add the docs in a list to bulk insert in the end.

Once we bulk insert sucessfully we update the publish time for next call which is going to be most recent doc

In this process we also maintain a LastInsertedIndex which is an auto incrementing variable that keeps track of the most recent event that
we have in our Database. This key is used to page results on the API.
*/

type WorkerHandler struct {
	query string

	apiKeys               []string
	apiKeyIndex           int
	currentPublishedTime  time.Time
	previousPublishedTime time.Time
	nextPageToken         string
	youtubeHandler        *youtube_handler.YoutubeHandler
	videoMetadataImpl     storage.VideoMetadataInterface
}

func NewWorkerHandler(apiKeys []string) (*WorkerHandler, error) {
	return &WorkerHandler{
		query:                "official|cricket|football|tennis|boating|sailing|food|minecraft|gaming|news|new",
		currentPublishedTime: time.Now().UTC().Add(-1 * time.Hour),
		apiKeys:              apiKeys,
		apiKeyIndex:          0,
		youtubeHandler:       youtube_handler.NewYoutubeHandler(apiKeys[0]),
		videoMetadataImpl:    storage.NewVideoMetadataImpl(),
	}, nil
}

// Starts the worker and executes it at an interval of 10 seconds
func (h *WorkerHandler) Start() {
	logger := common.GetLogger()

	for {
		err := h.Execute()
		if err != nil {
			if strings.Contains(err.Error(), "quotaExceeded") {
				logger.Info("API Key Quota Exceeded")
				h.youtubeHandler.UpdateAPIKey(h.FetchNextAPIKey())
			} else {
				logger.WithError(err).Error("Error in worker, exiting worker")
				break
			}
		}
		time.Sleep(10 * time.Second)
	}

}

// Updates the API key in case quota was exceeded in the current key
func (h *WorkerHandler) FetchNextAPIKey() string {
	logger := common.GetLogger()
	h.apiKeyIndex += 1
	if h.apiKeyIndex == len(h.apiKeys) {
		h.apiKeyIndex = 0
	}

	logger.WithField("NewAPIKeyIndex", h.apiKeyIndex).WithField("TotalKeys", len(h.apiKeys)).Info("API Key Updated")
	return h.apiKeys[h.apiKeyIndex]
}

func (h *WorkerHandler) Execute() error {
	logger := common.GetLogger()
	var err error
	var results *youtube.SearchListResponse

	logger.WithField("From", h.currentPublishedTime).WithField("NextPageToken", h.nextPageToken).Info("Fetching Youtube Data")

	if h.nextPageToken != "" {
		results, err = h.youtubeHandler.DoSearchListNextPage(h.query, []string{"snippet"}, "video", "date", h.previousPublishedTime.Format(time.RFC3339), h.nextPageToken, 50)
	} else {
		results, err = h.youtubeHandler.DoSearchList(h.query, []string{"snippet"}, "video", "date", h.currentPublishedTime.Format(time.RFC3339), 50)
	}

	if err != nil {
		return err
	}

	metadataToInsert := []*storage.VideoMetadata{}

	for _, result := range results.Items {
		videoMetadata, err := newVideoMetadata(result)
		if err != nil {
			return err
		}

		storageMetadata, err := h.videoMetadataImpl.FindOneMetadataWithVideoID(videoMetadata.VideoID)
		if err != nil {
			return err
		}

		if storageMetadata != nil {
			if storageMetadata.PublishedAt.Before(videoMetadata.PublishedAt) {
				h.videoMetadataImpl.UpdateOneMetadata(videoMetadata.VideoID, videoMetadata)
			}
		} else {
			metadataToInsert = append(metadataToInsert, videoMetadata)
		}

		publishedAt, err := time.Parse(time.RFC3339, result.Snippet.PublishedAt)
		if err != nil {
			return err
		}

		// Save the publishedAt time for the next call. This will be only used
		// if there is no next page token. If there is a next page token, previousPublishedAt
		// along with the token will be used.
		if h.currentPublishedTime.Before(publishedAt) {
			h.previousPublishedTime = h.currentPublishedTime
			h.currentPublishedTime = publishedAt
		}
	}

	h.nextPageToken = ""
	if results.NextPageToken != "" {
		h.nextPageToken = results.NextPageToken
	}

	if len(metadataToInsert) > 0 {
		logger.WithField("InsertDocCount", len(metadataToInsert)).Info("Publishing documents to database")

		err := h.videoMetadataImpl.BulkInsertMetadata(metadataToInsert)
		if err != nil {
			return err
		}
	}
	return nil
}

// Takes the result from youtube API and populates in our structure format
func newVideoMetadata(result *youtube.SearchResult) (*storage.VideoMetadata, error) {
	publishedAtTime, err := time.Parse(time.RFC3339, result.Snippet.PublishedAt)
	if err != nil {
		return nil, err
	}

	videoData := &storage.VideoMetadata{
		VideoID:     result.Id.VideoId,
		Title:       result.Snippet.Title,
		Description: result.Snippet.Description,
		PublishedAt: publishedAtTime,
	}

	if result.Snippet.Thumbnails.Default != nil {
		videoData.DefaultThumbnailURL = result.Snippet.Thumbnails.Default.Url
	}
	if result.Snippet.Thumbnails.Maxres != nil {
		videoData.MaxresThumbnailURL = result.Snippet.Thumbnails.Maxres.Url
	}
	if result.Snippet.Thumbnails.High != nil {
		videoData.HighThumbnailURL = result.Snippet.Thumbnails.High.Url
	}
	if result.Snippet.Thumbnails.Medium != nil {
		videoData.MediumThumbnailURL = result.Snippet.Thumbnails.Medium.Url
	}
	if result.Snippet.Thumbnails.Standard != nil {
		videoData.StandardThumbnailURL = result.Snippet.Thumbnails.Standard.Url
	}

	return videoData, nil
}
