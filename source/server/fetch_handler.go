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

type FetchResponse struct {
	User     string
	Page     int
	Metadata []*storage.VideoMetadata
}

func (h *ServerHandler) NewFetchHandler(w http.ResponseWriter, r *http.Request) {
	logger := common.GetLogger()
	config := common.GetConfiguration()
	var err error

	pagesizeParam := r.URL.Query().Get("pagesize")
	var pageSize int
	if pagesizeParam == "" {
		pageSize = config.DefaultPageSize
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
		// Update
	} else {
		h.userHandler.CreateUser(&storage.User{
			UserID:    userID,
			PageSize:  pageSize,
			Timestamp: time.Now().UTC(),
		})
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
