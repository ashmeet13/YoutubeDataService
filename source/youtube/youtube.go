package youtube

import (
	"context"

	"github.com/ashmeet13/YoutubeDataService/source/common"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

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

func (h *YoutubeHandler) DoSearchList(query string, parts []string, resourceType string, orderBy string, publishedAfter string) (*youtube.SearchListResponse, error) {
	searchRequest := h.youtubeClient.Search.List(parts).Q(query).
		Type(resourceType).Order(orderBy).PublishedAfter(publishedAfter)

	response, err := searchRequest.Do()

	if err != nil {
		return nil, err
	}

	return response, nil
}
