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

We query the GET /search/list endpoint to fetch data with the following parameters set
1. Parts = snippet, this fetches us the metadata we require
2. Type = video, we only want data for videos
3. OrderBy = date, we get the data ordered with most recent event at the start of the response
4. PublishedAfter = <time.RFC3339 formated date> we will be fetching all the data that was created on and after this point of time.
5. MaxResults = 50

Assumption - /search/list can return data for a video which has been updated

We iterate over the result and create structs to insert into our DB.
While doing so we make a DB read call to check if the video already exists.

If it does and the published at of the new result is after the database doc, we update the doc for the video id.
Else, we add the docs in a list to bulk insert in the end.

While iterating we also check if the results were for a call the nextPage of previous call or is it a fresh call
If it's a fresh call then we have recieved data ordered by date and we would use this data to set the next timestamp
of our call

Finally we check
*/

type WorkerHandler struct {
	query string

	apiKeys       []string
	nextPageToken string

	sleepTime   int
	apiKeyIndex int

	currentPublishedTime  time.Time
	previousPublishedTime time.Time

	youtubeHandler       youtube_handler.YoutubeInterface
	videoMetadataHandler storage.VideoMetadataInterface
}

func NewWorkerHandler(query string, apiKeys []string) (*WorkerHandler, error) {
	return &WorkerHandler{
		query:                query,
		currentPublishedTime: time.Now().UTC(),
		apiKeys:              apiKeys,
		apiKeyIndex:          0,
		sleepTime:            10,
		youtubeHandler:       youtube_handler.NewYoutubeHandler(apiKeys[0]),
		videoMetadataHandler: storage.NewVideoMetadataImpl(),
		nextPageToken:        "",
	}, nil
}

// Starts the worker and executes it at an interval of 10 seconds or 5 seconds
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
		logger.WithField("SleepDuration", h.sleepTime).Info("Worker Execution Completed")
		time.Sleep(time.Duration(h.sleepTime) * time.Second)
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

// Executes - To Fetch Data and Publish to DB
func (h *WorkerHandler) Execute() error {
	// Reset sleep time for next call
	h.sleepTime = 10

	logger := common.GetLogger()

	var err error
	var results *youtube.SearchListResponse

	// 1. If NextPageToken present
	//			Request Paginated Response from PreviousPublishedAt, i.e. previous call next page
	//    Else,
	// 			Request Response from CurrentPublishedAt, i.e. fresh call
	if h.nextPageToken != "" {
		logger.WithField("From", h.currentPublishedTime).WithField("NextPageToken", h.nextPageToken).Info("Fetching Youtube Data for next Page")
		results, err = h.youtubeHandler.DoSearchListNextPage(h.query, []string{"snippet"}, "video", "date", h.previousPublishedTime.Format(time.RFC3339), h.nextPageToken, 50)
	} else {
		logger.WithField("From", h.currentPublishedTime).Info("Fetching Youtube Data for new DateTime")
		results, err = h.youtubeHandler.DoSearchList(h.query, []string{"snippet"}, "video", "date", h.currentPublishedTime.Format(time.RFC3339), 50)
	}

	if err != nil {
		return err
	}

	logger.Info("Recieved Youtube Result")

	metadataToInsert := []*storage.VideoMetadata{}
	// 2. Iterate over the results.Items to fetch required Data
	for _, result := range results.Items {

		// 3. Format it into required Struct
		videoMetadata, err := newVideoMetadata(result)
		if err != nil {
			return err
		}

		// 4. Check if video already present
		storageMetadata, err := h.videoMetadataHandler.FindOneMetadataWithVideoID(videoMetadata.VideoID)
		if err != nil {
			return err
		}

		if storageMetadata != nil {
			// 5. If present and has been updated, update value in DB
			if storageMetadata.PublishedAt.Before(videoMetadata.PublishedAt) {
				logger.WithField("VideoID", videoMetadata.VideoID).Info("Updating Document")
				h.videoMetadataHandler.UpdateOneMetadata(videoMetadata.VideoID, videoMetadata)
			}
		} else {
			// 6. If not present, add in list to bulk insert later
			metadataToInsert = append(metadataToInsert, videoMetadata)
		}

		// 7. If no nextPageToken was used, update the currentPublishedTime to use for next call fresh call
		if h.nextPageToken == "" {
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
	}

	// 8. Reset token for next call, if a page token is available set it and reduce sleep time
	h.nextPageToken = ""
	if results.NextPageToken != "" {
		h.nextPageToken = results.NextPageToken
		h.sleepTime = 5
	}

	// 9. Bulk Insert into DB
	if len(metadataToInsert) > 0 {
		logger.WithField("InsertDocumentCount", len(metadataToInsert)).Info("Publishing documents to database")

		err := h.videoMetadataHandler.BulkInsertMetadata(metadataToInsert)
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

	if result.Snippet.Thumbnails != nil {
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
	}

	return videoData, nil
}
