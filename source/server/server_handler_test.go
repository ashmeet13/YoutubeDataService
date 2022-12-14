package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ashmeet13/YoutubeDataService/source/common"
	"github.com/ashmeet13/YoutubeDataService/source/storage"
	"github.com/ashmeet13/YoutubeDataService/source/storage/mock_storage"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ServerHandlerSuite struct {
	suite.Suite
	*require.Assertions
	ctrl *gomock.Controller

	mockVideoMetadataStore *mock_storage.MockVideoMetadataInterface
	mockUserStore          *mock_storage.MockUserInterface
	serverHandler          *ServerHandler
}

func TestWorkerHandlerSuite(t *testing.T) {
	suite.Run(t, new(ServerHandlerSuite))
}

func (s *ServerHandlerSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.ctrl = gomock.NewController(s.T())

	s.mockUserStore = mock_storage.NewMockUserInterface(s.ctrl)
	s.mockVideoMetadataStore = mock_storage.NewMockVideoMetadataInterface(s.ctrl)

	s.serverHandler = &ServerHandler{
		userHandler:          s.mockUserStore,
		videoMetadataHandler: s.mockVideoMetadataStore,

		config: &common.Configuration{
			DefaultPageSize: 25,
		},
	}
}

func (s *ServerHandlerSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func (s *ServerHandlerSuite) TestSearchHandler_WrongContentType() {
	req := httptest.NewRequest(http.MethodPost, "/search", nil)
	res := httptest.NewRecorder()

	s.serverHandler.SearchHandler(res, req)

	message, _ := ioutil.ReadAll(res.Body)

	s.Equal(http.StatusUnsupportedMediaType, res.Code)
	s.Equal("Content-Type header is not application/json\n", string(message))
}

func (s *ServerHandlerSuite) TestSearchHandler_NoBody() {
	req := httptest.NewRequest(http.MethodPost, "/search", nil)
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	s.serverHandler.SearchHandler(res, req)

	message, _ := ioutil.ReadAll(res.Body)

	s.Equal(http.StatusBadRequest, res.Code)
	s.Equal("Failed to read request body\n", string(message))
}

func (s *ServerHandlerSuite) TestSearchHandler_EmptyBody() {
	testSearchFilters := &SearchFilters{
		Title:       "",
		Description: "",
	}

	jsonRequest, _ := json.Marshal(testSearchFilters)

	req := httptest.NewRequest(http.MethodPost, "/search", bytes.NewBuffer(jsonRequest))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	s.serverHandler.SearchHandler(res, req)

	message, _ := ioutil.ReadAll(res.Body)

	s.Equal(http.StatusBadRequest, res.Code)
	s.Equal("Title and Description both cannot be empty\n", string(message))
}

func (s *ServerHandlerSuite) TestSearchHandler_DBFail() {
	testSearchFilters := &SearchFilters{
		Title:       "test_title",
		Description: "test_description",
	}

	jsonRequest, _ := json.Marshal(testSearchFilters)

	req := httptest.NewRequest(http.MethodPost, "/search", bytes.NewBuffer(jsonRequest))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	s.mockVideoMetadataStore.EXPECT().FindMetadataTextSearch("test_title").Return(nil, errors.New("dummy test error"))
	s.serverHandler.SearchHandler(res, req)

	message, _ := ioutil.ReadAll(res.Body)

	s.Equal(http.StatusInternalServerError, res.Code)
	s.Equal("dummy test error\n", string(message))
}

func (s *ServerHandlerSuite) TestSearchHandler_Ok() {
	testSearchFilters := &SearchFilters{
		Title:       "test_title",
		Description: "test_description",
	}

	jsonRequest, _ := json.Marshal(testSearchFilters)

	req := httptest.NewRequest(http.MethodPost, "/search", bytes.NewBuffer(jsonRequest))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	returnMetadata1 := &storage.VideoMetadata{
		Title:       "test_title1",
		Description: "test_description1",
		VideoID:     "123",
	}

	returnMetadata2 := &storage.VideoMetadata{
		Title:       "test_title2",
		Description: "test_description2",
		VideoID:     "456",
	}

	s.mockVideoMetadataStore.EXPECT().FindMetadataTextSearch("test_title").Return([]*storage.VideoMetadata{returnMetadata1}, nil)
	s.mockVideoMetadataStore.EXPECT().FindMetadataTextSearch("test_description").Return([]*storage.VideoMetadata{returnMetadata2}, nil)
	s.serverHandler.SearchHandler(res, req)

	var response SearchResponse
	_ = json.NewDecoder(res.Body).Decode(&response)

	s.Equal(http.StatusOK, res.Code)
	s.Equal(2, len(response.Metadata))
	s.Equal("test_title1", response.Metadata[0].Title)
	s.Equal("test_description1", response.Metadata[0].Description)
	s.Equal("123", response.Metadata[0].VideoID)

	s.Equal("test_title2", response.Metadata[1].Title)
	s.Equal("test_description2", response.Metadata[1].Description)
	s.Equal("456", response.Metadata[1].VideoID)
}

func (s *ServerHandlerSuite) TestNewFetchHandler_Ok() {
	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1/fetch?userid=12345&pagesize=15", nil)
	res := httptest.NewRecorder()

	// Mocking new user behaviour
	s.mockUserStore.EXPECT().ReadUser("12345").Return(nil, nil)
	s.mockUserStore.EXPECT().CreateUser(gomock.Any()).Return(nil)

	s.serverHandler.NewFetchHandler(res, req)

	var response FetchResponse
	_ = json.NewDecoder(res.Body).Decode(&response)

	s.Equal(http.StatusOK, res.Code)
	s.Equal(response.User, "12345")
	s.Equal(response.Page, 0)
}

func (s *ServerHandlerSuite) TestNewFetchHandler_WrongPageSize() {
	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1/fetch?userid=12345&pagesize=ab", nil)
	res := httptest.NewRecorder()

	s.serverHandler.NewFetchHandler(res, req)

	s.Equal(http.StatusInternalServerError, res.Code)
}

func (s *ServerHandlerSuite) TestNewFetchHandler_OK_ExistingUser() {
	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1/fetch?userid=12345", nil)
	res := httptest.NewRecorder()

	s.mockUserStore.EXPECT().ReadUser("12345").Return(&storage.User{UserID: "12345"}, nil)
	s.mockUserStore.EXPECT().UpdateUser("12345", gomock.Any()).Return(nil)

	s.serverHandler.NewFetchHandler(res, req)

	var response FetchResponse
	_ = json.NewDecoder(res.Body).Decode(&response)

	s.Equal(http.StatusOK, res.Code)
	s.Equal(response.User, "12345")
	s.Equal(response.Page, 0)
}

func (s *ServerHandlerSuite) TestFetchHandler_OK() {
	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1/fetch/12345/2", nil)
	res := httptest.NewRecorder()

	vars := map[string]string{
		"userid": "12345",
		"page":   "2",
	}

	req = mux.SetURLVars(req, vars)

	user := &storage.User{
		UserID:    "12345",
		PageSize:  5,
		Timestamp: time.Now().UTC(),
	}

	metadata := []*storage.VideoMetadata{
		{
			VideoID: "abc",
		},
		{
			VideoID: "def",
		},
	}

	s.mockUserStore.EXPECT().ReadUser("12345").Return(user, nil)
	s.mockVideoMetadataStore.EXPECT().FetchPagedMetadata(user.Timestamp, int64(5), int64(5)).Return(metadata, nil)

	s.serverHandler.FetchHandler(res, req)

	var response FetchResponse
	_ = json.NewDecoder(res.Body).Decode(&response)

	s.Equal(http.StatusOK, res.Code)
	s.Equal("12345", response.User)
	s.Equal(2, response.Page)
	s.Equal(2, len(response.Metadata))
	s.Equal("abc", response.Metadata[0].VideoID)
	s.Equal("def", response.Metadata[1].VideoID)
}
