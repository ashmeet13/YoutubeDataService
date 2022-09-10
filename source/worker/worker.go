package worker

import (
	"fmt"
	"strings"
	"time"

	"github.com/ashmeet13/YoutubeDataService/source/common"
	"github.com/ashmeet13/YoutubeDataService/source/storage"
	"github.com/ashmeet13/YoutubeDataService/source/youtube"
)

type WorkerHandler struct {
	query          string
	publishedTime  string
	apiKeys        []string
	apiKeyIndex    int
	youtubeHandler *youtube.YoutubeHandler
}

func NewWorkerHandler(apiKeys []string) *WorkerHandler {
	fmt.Println(apiKeys)
	return &WorkerHandler{
		query:          "official",
		publishedTime:  time.Now().UTC().Format(time.RFC3339),
		apiKeys:        apiKeys,
		apiKeyIndex:    0,
		youtubeHandler: youtube.NewYoutubeHandler(apiKeys[0]),
	}
}

func (h *WorkerHandler) Start() {
	logger := common.GetLogger()

	for {
		err := h.Execute()

		if err != nil {
			if strings.Contains(err.Error(), "quotaExceeded") {
				logger.Info("API Key Quota Exceeded")
				h.youtubeHandler.UpdateAPIKey(h.FetchNextAPIKey())
				logger.WithField("ApiKeyIndex", h.apiKeyIndex).Info("API Key Updated")
			} else {
				logger.WithError(err).Error("Exiting Worker")
				break
			}
		}
		time.Sleep(time.Second * 10)
	}
}

func (h *WorkerHandler) FetchNextAPIKey() string {
	h.apiKeyIndex += 1
	if h.apiKeyIndex == len(h.apiKeys) {
		h.apiKeyIndex = 0
	}

	return h.apiKeys[h.apiKeyIndex]
}

func (h *WorkerHandler) Execute() error {
	logger := common.GetLogger()
	logger = logger.WithField("function", "Worker/Execute")
	count := 0

	logger.WithField("PublishedAt", h.publishedTime).Info("Fetching Youtube Data")
	results, err := h.youtubeHandler.DoSearchList(h.query, []string{"snippet"}, "video", "date", h.publishedTime)
	if err != nil {
		return err
	}

	for _, result := range results.Items {
		videoData := &storage.VideoMetadata{
			Title:       result.Snippet.Title,
			Description: result.Snippet.Description,
			PublishedAt: result.Snippet.PublishedAt,
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
			videoData.StandarThumbnailURL = result.Snippet.Thumbnails.Standard.Url
		}

		_, err := storage.InsertOne(storage.VideoMetadataC, videoData)
		if err != nil {
			logger.Error("Failed to insert document", result.Id)
			return err
		} else {
			h.publishedTime = result.Snippet.PublishedAt
			count += 1
		}
	}

	logger.WithField("InsertedDocs", count).WithField("Page", results.NextPageToken).Info("Docs Inserted")

	return nil
}
