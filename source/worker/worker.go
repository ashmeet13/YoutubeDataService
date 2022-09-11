package worker

import (
	"strings"
	"time"

	"github.com/ashmeet13/YoutubeDataService/source/common"
	"github.com/ashmeet13/YoutubeDataService/source/storage"
	"github.com/ashmeet13/YoutubeDataService/source/youtube"
)

type WorkerHandler struct {
	query             string
	publishedTime     string
	apiKeys           []string
	apiKeyIndex       int
	lastInsertedIndex int
	youtubeHandler    *youtube.YoutubeHandler
	videoMetadataImpl storage.VideoMetadataInterface
}

func NewWorkerHandler(apiKeys []string) (*WorkerHandler, error) {
	videoMetadataImpl := storage.NewVideoMetadataImpl()

	lastInsertedIndex, err := videoMetadataImpl.FindLastInsertedIndex()
	if err != nil {
		return nil, err
	}

	return &WorkerHandler{
		query:             "official|cricket|football|tennis|boating|sailing|food|minecraft|gaming|news|new",
		publishedTime:     time.Now().UTC().Format(time.RFC3339),
		apiKeys:           apiKeys,
		apiKeyIndex:       0,
		lastInsertedIndex: lastInsertedIndex,
		youtubeHandler:    youtube.NewYoutubeHandler(apiKeys[0]),
		videoMetadataImpl: storage.NewVideoMetadataImpl(),
	}, nil
}

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
	logger = logger.WithField("function", "Worker/Execute")

	var err error
	videoMetadatas := []*storage.VideoMetadata{}

	logger.WithField("From", h.publishedTime).Info("Fetching Youtube Data")
	results, err := h.youtubeHandler.DoSearchList(h.query, []string{"snippet"}, "video", "date", h.publishedTime)
	if err != nil {
		return err
	}

	lastIndex := len(results.Items) - 1
	for index := range results.Items {
		result := results.Items[lastIndex-index]
		publishedAtTime, err := time.Parse(time.RFC3339, result.Snippet.PublishedAt)
		if err != nil {
			return err
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

		videoMetadatas = append(videoMetadatas, videoData)
	}

	if len(videoMetadatas) > 0 {
		h.publishedTime = results.Items[0].Snippet.PublishedAt
		logger.WithField("PublishedTime", h.publishedTime).Info("Updated Published Time")
		err = h.PublishToDatabase(videoMetadatas)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *WorkerHandler) PublishToDatabase(videoMetadatas []*storage.VideoMetadata) error {
	logger := common.GetLogger()
	logger.WithField("TotalDocCount", len(videoMetadatas)).Info("New Publish Request")

	toInsert := []*storage.VideoMetadata{}
	for _, metadata := range videoMetadatas {
		exists, err := h.videoMetadataImpl.MetadataExists(metadata.VideoID)
		if err != nil {
			return err
		}
		if !exists {
			metadata.DocumentIndex = h.lastInsertedIndex + 1
			h.lastInsertedIndex = metadata.DocumentIndex
			toInsert = append(toInsert, metadata)
		}
	}

	if len(toInsert) > 0 {
		logger.WithField("InsertDocCount", len(toInsert)).WithField("LastIndex", h.lastInsertedIndex).Info("Publishing documents to database")
		err := h.videoMetadataImpl.BulkInsertMetadata(toInsert)
		if err != nil {
			return err
		}
	}
	return nil
}
