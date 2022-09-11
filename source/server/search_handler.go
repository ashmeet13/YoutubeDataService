package server

import (
	"encoding/json"
	"net/http"

	"github.com/ashmeet13/YoutubeDataService/source/common"
	"github.com/ashmeet13/YoutubeDataService/source/storage"
)

type SearchFilters struct {
	Title       string
	Description string
}

type SearchResponse struct {
	Message  string
	Metadata []*storage.VideoMetadata
}

func (h *ServerHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	logger := common.GetLogger()

	if r.Header.Get("Content-Type") == "" || r.Header.Get("Content-Type") != "application/json" {
		msg := "Content-Type header is not application/json"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	var searchFilters SearchFilters

	err := json.NewDecoder(r.Body).Decode(&searchFilters)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if searchFilters.Title == "" && searchFilters.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Title and Description both cannot be empty", http.StatusBadRequest)
		return
	}

	metadata, err := h.storageHandler.FindOneMetadata(searchFilters.Title, searchFilters.Description)
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
