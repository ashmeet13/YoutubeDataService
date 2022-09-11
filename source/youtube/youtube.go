package youtube_handler

import (
	"context"

	"github.com/ashmeet13/YoutubeDataService/source/common"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeHandler struct {
	apiKey        string
	youtubeClient *youtube.Service
}

func NewYoutubeHandler(apiKey string) *YoutubeHandler {
	youtubeClient, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))

	if err != nil {
		panic("Failed to create youtube service client")
	}

	return &YoutubeHandler{
		apiKey:        apiKey,
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
	h.apiKey = apiKey
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

// // https://youtube.googleapis.com/youtube/v3/search?part=snippet&order=date&publishedAfter=2022-09-10T18%3A21%3A38Z&q=official&type=video&key=[YOUR_API_KEY]
// func (h *YoutubeHandler) DoSearchListNew(query string, parts []string, resourceType string, orderBy string, publishedAfter string) {
// 	URL := fmt.Sprintf(
// 		"https://youtube.googleapis.com/youtube/v3/search?part=%s&order=%s&publishedAfter=%s&q=%s&type=%s&key=%s",
// 		strings.Join(parts, ","),
// 		orderBy,
// 		publishedAfter,
// 		query,
// 		resourceType,
// 		h.apiKey,
// 	)

// 	fmt.Println(URL)

// 	resp, err := http.Get(URL)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(string(body))

// 	// searchRequest := h.youtubeClient.Search.List(parts).Q(query).
// 	// 	Type(resourceType).Order(orderBy).PublishedAfter(publishedAfter)

// 	// response, err := searchRequest.Do()

// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// return response, nil
// }
