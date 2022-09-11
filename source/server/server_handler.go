package server

import (
	"encoding/json"
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
	Message  string
	Metadata []*storage.VideoMetadata
}

type FetchResponse struct {
	User     string
	Page     int
	Metadata []*storage.VideoMetadata
}

func (h *ServerHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	logger := common.GetLogger()

	if r.Header.Get("Content-Type") == "" || r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		msg := "Content-Type header is not application/json"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	var searchFilters SearchFilters

	err := json.NewDecoder(r.Body).Decode(&searchFilters)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := "Failed to read request body"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if searchFilters.Title == "" && searchFilters.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Title and Description both cannot be empty", http.StatusBadRequest)
		return
	}

	metadata, err := h.videoMetadataHandler.FindOneMetadata(searchFilters.Title, searchFilters.Description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := &SearchResponse{
		Message: "Ok",
	}

	if metadata != nil {
		response.Metadata = []*storage.VideoMetadata{metadata}
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		logger.WithError(err).Error("Error happened in JSON marshal")
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
