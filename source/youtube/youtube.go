package youtube_handler

import (
	"context"

	"github.com/ashmeet13/YoutubeDataService/source/common"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

//go:generate mockgen --destination=./mock_youtube/youtube.go github.com/ashmeet13/YoutubeDataService/source/youtube YoutubeInterface
type YoutubeInterface interface {
	UpdateAPIKey(apiKey string) error
	DoSearchList(query string, parts []string, resourceType string, orderBy string, publishedAfter string, maxResults int) (*youtube.SearchListResponse, error)
	DoSearchListNextPage(query string, parts []string, resourceType string, orderBy string, publishedAfter string, nextPageToken string, maxResults int) (*youtube.SearchListResponse, error)
}

type YoutubeHandler struct {
	youtubeClient *youtube.Service
}

func NewYoutubeHandler(apiKey string) *YoutubeHandler {
	youtubeClient, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))

	if err != nil {
		panic("Failed to create youtube service client")
	}

	return &YoutubeHandler{
		youtubeClient: youtubeClient,
	}
}

func (h *YoutubeHandler) UpdateAPIKey(apiKey string) error {
	logger := common.GetLogger()
	youtubeClient, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))

	if err != nil {
		logger.WithError(err).Error("failed to create youtube service client")
	}

	h.youtubeClient = youtubeClient
	return nil
}

func (h *YoutubeHandler) DoSearchList(query string, parts []string, resourceType string, orderBy string, publishedAfter string, maxResults int) (*youtube.SearchListResponse, error) {
	searchRequest := h.youtubeClient.Search.List(parts).Q(query).
		Type(resourceType).Order(orderBy).PublishedAfter(publishedAfter).MaxResults(int64(maxResults))

	response, err := searchRequest.Do()

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (h *YoutubeHandler) DoSearchListNextPage(query string, parts []string, resourceType string, orderBy string, publishedAfter string, nextPageToken string, maxResults int) (*youtube.SearchListResponse, error) {
	searchRequest := h.youtubeClient.Search.List(parts).Q(query).
		Type(resourceType).Order(orderBy).PublishedAfter(publishedAfter).MaxResults(int64(maxResults)).PageToken(nextPageToken)

	response, err := searchRequest.Do()

	if err != nil {
		return nil, err
	}

	return response, nil
}
