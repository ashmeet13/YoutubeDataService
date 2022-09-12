package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ashmeet13/YoutubeDataService/source/common"
	"github.com/ashmeet13/YoutubeDataService/source/storage"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func NewServerHandler() *ServerHandler {
	return &ServerHandler{
		config: common.SetupConfiguration(),

		videoMetadataHandler: storage.NewVideoMetadataImpl(),
		userHandler:          storage.NewUserImpl(),
	}
}

type ServerHandler struct {
	config               *common.Configuration
	videoMetadataHandler storage.VideoMetadataInterface
	userHandler          storage.UserInterface
}

type SearchFilters struct {
	Title       string
	Description string
}

type SearchResponse struct {
	Metadata []*storage.VideoMetadata
}

type FetchResponse struct {
	User     string
	Page     int
	Metadata []*storage.VideoMetadata
}

// Handles /search request
func (h *ServerHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	logger := common.GetLogger()
	logger.Info("New Search Request")

	var err error
	var matchedDocs []*storage.VideoMetadata

	// Make sure JSON body is present
	if r.Header.Get("Content-Type") == "" || r.Header.Get("Content-Type") != "application/json" {
		logger.Error("Content Type not correct")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		msg := "Content-Type header is not application/json"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	var searchFilters SearchFilters

	// Fetch and decode the JSON Body
	err = json.NewDecoder(r.Body).Decode(&searchFilters)
	if err != nil {
		logger.WithError(err).Error("Failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		msg := "Failed to read request body"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// Make sure we have some search parameter
	if searchFilters.Title == "" && searchFilters.Description == "" {
		logger.Error("Both title and description are empty")
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Title and Description both cannot be empty", http.StatusBadRequest)
		return
	}

	logger = logger.WithField("Title", searchFilters.Title).
		WithField("Description", searchFilters.Description)

	// Make DB call to search matching titles
	if searchFilters.Title != "" {
		logger.Info("Searching data matching the title")
		titleMatchedDocs, err := h.videoMetadataHandler.FindMetadataTextSearch(searchFilters.Title)
		if err != nil {
			logger.WithError(err).Error("Failed to get data from database")
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		matchedDocs = append(matchedDocs, titleMatchedDocs...)
	}

	// Make DB call to search matching descriptions
	if searchFilters.Description != "" {
		logger.Info("Searching data matching the description")
		desMatchedDocs, err := h.videoMetadataHandler.FindMetadataTextSearch(searchFilters.Description)
		if err != nil {
			logger.WithError(err).Error("Failed to get data from database")
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		matchedDocs = append(matchedDocs, desMatchedDocs...)
	}

	// Build and return response
	response := &SearchResponse{
		Metadata: matchedDocs,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		logger.WithError(err).Error("Error happened in JSON marshal")
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (h *ServerHandler) NewFetchHandler(w http.ResponseWriter, r *http.Request) {
	logger := common.GetLogger()
	var err error

	pagesizeParam := r.URL.Query().Get("pagesize")
	var pageSize int
	if pagesizeParam == "" {
		pageSize = h.config.DefaultPageSize
	} else {
		pageSize, err = strconv.Atoi(pagesizeParam)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	userID := r.URL.Query().Get("userid")
	if userID == "" {
		userID = uuid.NewString()
	}

	user, err := h.userHandler.ReadUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user != nil {
		logger.WithField("UserID", userID).Info("Updating User")
		user.Timestamp = time.Now().UTC()
		err = h.userHandler.UpdateUser(userID, user)
	} else {
		logger.WithField("UserID", userID).Info("Creating User")
		err = h.userHandler.CreateUser(&storage.User{
			UserID:    userID,
			PageSize:  pageSize,
			Timestamp: time.Now().UTC(),
		})
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(&FetchResponse{
		User: userID,
		Page: 0,
	})

	if err != nil {
		logger.WithError(err).Error("Error happened in JSON marshal")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (h *ServerHandler) FetchHandler(w http.ResponseWriter, r *http.Request) {
	logger := common.GetLogger()

	vars := mux.Vars(r)
	userID, ok := vars["userid"]
	if !ok {
		msg := "userid is missing in parameters"
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	pageParam, ok := vars["page"]
	if !ok {
		msg := "page is missing in parameters"
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.WithField("User", userID).WithField("Page", page).Info("Fetch Request")

	user, err := h.userHandler.ReadUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Could not find user with userid %s", userID)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	offset := user.PageSize * (page - 1)

	metadata, err := h.videoMetadataHandler.FetchPagedMetadata(user.Timestamp, int64(offset), int64(user.PageSize))
	if err != nil {
		msg := "Failed in fetching page"
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(&FetchResponse{
		User:     userID,
		Page:     page,
		Metadata: metadata,
	})
	if err != nil {
		logger.WithError(err).Error("Error happened in JSON marshal")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
