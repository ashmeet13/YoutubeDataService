package worker

import (
	"testing"
	"time"

	"github.com/ashmeet13/YoutubeDataService/source/storage"
	"github.com/ashmeet13/YoutubeDataService/source/storage/mock_storage"
	"github.com/ashmeet13/YoutubeDataService/source/youtube/mock_youtube"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/api/youtube/v3"
)

type WorkerHandlerSuite struct {
	suite.Suite
	*require.Assertions
	ctrl *gomock.Controller

	mockVideoMetadataStore *mock_storage.MockVideoMetadataInterface
	mockYoutubeHandler     *mock_youtube.MockYoutubeInterface

	workerHandler *WorkerHandler
}

func TestWorkerHandlerSuite(t *testing.T) {
	suite.Run(t, new(WorkerHandlerSuite))
}

func (s *WorkerHandlerSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.ctrl = gomock.NewController(s.T())

	s.mockVideoMetadataStore = mock_storage.NewMockVideoMetadataInterface(s.ctrl)
	s.mockYoutubeHandler = mock_youtube.NewMockYoutubeInterface(s.ctrl)

	s.workerHandler = &WorkerHandler{
		videoMetadataHandler: s.mockVideoMetadataStore,
		youtubeHandler:       s.mockYoutubeHandler,

		apiKeys:     []string{"abcd", "edfg"},
		apiKeyIndex: 0,

		query: "query",
	}
}

func (s *WorkerHandlerSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func (s *WorkerHandlerSuite) TestFetchNextAPIKey() {
	apiKey := s.workerHandler.FetchNextAPIKey()

	s.Equal(apiKey, "edfg")
	s.Equal(s.workerHandler.apiKeyIndex, 1)

	apiKey = s.workerHandler.FetchNextAPIKey()

	s.Equal(apiKey, "abcd")
	s.Equal(s.workerHandler.apiKeyIndex, 0)
}

func (s *WorkerHandlerSuite) TestExecute_FreshCall() {
	s.workerHandler.nextPageToken = ""

	currentPublishedTime := time.Now().UTC()
	s.workerHandler.currentPublishedTime = currentPublishedTime

	expectedDate := s.workerHandler.currentPublishedTime.Format(time.RFC3339)

	testPublishedAtTime := s.workerHandler.currentPublishedTime.Add(5 * time.Second)

	results := &youtube.SearchListResponse{
		Items: []*youtube.SearchResult{
			{
				Id: &youtube.ResourceId{
					VideoId: "test_video_id",
				},
				Snippet: &youtube.SearchResultSnippet{
					Title:       "test_title",
					Description: "test_description",
					PublishedAt: testPublishedAtTime.Format(time.RFC3339),
					Thumbnails: &youtube.ThumbnailDetails{
						High: &youtube.Thumbnail{
							Url: "test_high_url",
						},
					},
				},
			},
		},
	}

	expectedNewDate, _ := time.Parse(time.RFC3339, testPublishedAtTime.Format(time.RFC3339))

	metadataVideo := &storage.VideoMetadata{
		VideoID:          "test_video_id",
		Title:            "test_title",
		Description:      "test_description",
		PublishedAt:      expectedNewDate,
		HighThumbnailURL: "test_high_url",
	}

	s.mockYoutubeHandler.EXPECT().DoSearchList("query", []string{"snippet"}, "video", "date", expectedDate, 50).Return(results, nil)
	s.mockVideoMetadataStore.EXPECT().FindOneMetadataWithVideoID("test_video_id").Return(nil, nil)
	s.mockVideoMetadataStore.EXPECT().BulkInsertMetadata([]*storage.VideoMetadata{metadataVideo}).Return(nil)

	s.workerHandler.Execute()

	s.Equal("", s.workerHandler.nextPageToken)
	s.Equal(10, s.workerHandler.sleepTime)
	s.Equal(expectedNewDate, s.workerHandler.currentPublishedTime)
	s.Equal(currentPublishedTime, s.workerHandler.previousPublishedTime)
}

func (s *WorkerHandlerSuite) TestExecute_PageCall() {
	s.workerHandler.nextPageToken = "ABCD"

	currentPublishedTime := time.Now().UTC()
	s.workerHandler.currentPublishedTime = currentPublishedTime

	previousPublishedTime := currentPublishedTime.Add(-5 * time.Second)
	s.workerHandler.previousPublishedTime = previousPublishedTime

	expectedDate := s.workerHandler.previousPublishedTime.Format(time.RFC3339)

	testPublishedAtTime := s.workerHandler.currentPublishedTime.Add(5 * time.Second)

	results := &youtube.SearchListResponse{
		Items: []*youtube.SearchResult{
			{
				Id: &youtube.ResourceId{
					VideoId: "test_video_id",
				},
				Snippet: &youtube.SearchResultSnippet{
					Title:       "test_title",
					Description: "test_description",
					PublishedAt: testPublishedAtTime.Format(time.RFC3339),
					Thumbnails: &youtube.ThumbnailDetails{
						High: &youtube.Thumbnail{
							Url: "test_high_url",
						},
					},
				},
			},
		},
	}

	expectedNewDate, _ := time.Parse(time.RFC3339, testPublishedAtTime.Format(time.RFC3339))

	metadataVideo := &storage.VideoMetadata{
		VideoID:          "test_video_id",
		Title:            "test_title",
		Description:      "test_description",
		PublishedAt:      expectedNewDate,
		HighThumbnailURL: "test_high_url",
	}

	s.mockYoutubeHandler.EXPECT().DoSearchListNextPage("query", []string{"snippet"}, "video", "date", expectedDate, "ABCD", 50).Return(results, nil)
	s.mockVideoMetadataStore.EXPECT().FindOneMetadataWithVideoID("test_video_id").Return(nil, nil)
	s.mockVideoMetadataStore.EXPECT().BulkInsertMetadata([]*storage.VideoMetadata{metadataVideo}).Return(nil)

	s.workerHandler.Execute()

	s.Equal("", s.workerHandler.nextPageToken)
	s.Equal(10, s.workerHandler.sleepTime)
	s.Equal(currentPublishedTime, s.workerHandler.currentPublishedTime)
	s.Equal(previousPublishedTime, s.workerHandler.previousPublishedTime)
}
