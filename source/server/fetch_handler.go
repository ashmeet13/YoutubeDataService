package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	var err error

	pagesizeParam := r.URL.Query().Get("pagesize")
	var pagesize int
	if pagesizeParam == "" {
		pagesize = 2
	} else {
		pagesize, err = strconv.Atoi(pagesizeParam)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	userID := uuid.NewString()

	lastIndex, err := h.storageHandler.FindLastInsertedIndex()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.lock.Lock()
	defer h.lock.Unlock()

	h.userLastSentIndex[userID] = lastIndex
	h.userPageSize[userID] = pagesize

	for key, element := range h.userLastSentIndex {
		fmt.Println("Key:", key, "=>", "Element:", element)
	}
	for key, element := range h.userPageSize {
		fmt.Println("Key:", key, "=>", "Element:", element)
	}

	response := &FetchResponse{
		User: userID,
		Page: 0,
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		logger.WithError(err).Error("Error happened in JSON marshal")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(jsonResponse)
}

func (h *ServerHandler) FetchHandler(w http.ResponseWriter, r *http.Request) {
	for key, element := range h.userLastSentIndex {
		fmt.Println("Key:", key, "=>", "Element:", element)
	}
	for key, element := range h.userPageSize {
		fmt.Println("Key:", key, "=>", "Element:", element)
	}
	logger := common.GetLogger()

	vars := mux.Vars(r)
	userID, ok := vars["userid"]
	if !ok {
		msg := "userid is missing in parameters"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	page, ok := vars["page"]
	if !ok {
		msg := "page is missing in parameters"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	logger.WithField("User", userID).WithField("Page", page).Info("Fetch Request")

	lastIndex, ok := h.userLastSentIndex[userID]
	if !ok {
		msg := "No last index found"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	pageSize, ok := h.userPageSize[userID]
	if !ok {
		msg := "No page size found"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	fmt.Println(lastIndex, pageSize)

	w.WriteHeader(http.StatusOK)

}
