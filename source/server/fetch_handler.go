package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ashmeet13/YoutubeDataService/source/common"
	"github.com/ashmeet13/YoutubeDataService/source/storage"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type FetchResponse struct {
	User        string
	Page        int
	LatestIndex int
	Metadata    []*storage.VideoMetadata
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

	lastInsertedMetadata, err := h.storageHandler.FindLastInsertedMetadata()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastIndex := 0
	if lastInsertedMetadata != nil {
		lastIndex = lastInsertedMetadata.DocumentIndex
	}

	h.setUserPageSize(userID, pageSize)
	h.setUserEndIndex(userID, lastIndex)

	jsonResponse, err := json.Marshal(&FetchResponse{
		User:        userID,
		Page:        0,
		LatestIndex: lastIndex,
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

	userEndIndex, err := h.getUserEndIndex(userID)
	if err != nil {
		msg := "No starting index found"
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	pageSize, err := h.getUserPageSize(userID)
	if err != nil {
		msg := "No page size found"
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	requiredEndIndex := userEndIndex - (pageSize * (page - 1))
	requiredStartIndex := requiredEndIndex - pageSize

	metadata, err := h.storageHandler.FetchPage(requiredStartIndex, requiredEndIndex)
	if err != nil {
		msg := "Failed in fetching page"
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(&FetchResponse{
		User:        userID,
		Page:        page,
		Metadata:    metadata,
		LatestIndex: requiredEndIndex,
	})
	if err != nil {
		logger.WithError(err).Error("Error happened in JSON marshal")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (h *ServerHandler) setUserEndIndex(userID string, index int) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.userEndingIndex[userID] = index
}

func (h *ServerHandler) getUserEndIndex(userID string) (int, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	value, ok := h.userEndingIndex[userID]
	if !ok {
		return 0, errors.New("no starting index found")
	}

	return value, nil
}

func (h *ServerHandler) setUserPageSize(userID string, pageSize int) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.userPageSize[userID] = pageSize
}

func (h *ServerHandler) getUserPageSize(userID string) (int, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	value, ok := h.userPageSize[userID]
	if !ok {
		return 0, errors.New("no starting index found")
	}

	return value, nil
}
