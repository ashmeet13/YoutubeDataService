package server

import (
	"net/http"

	"github.com/ashmeet13/YoutubeDataService/source/common"
	"github.com/gorilla/mux"
)

func Start() {
	logger := common.GetLogger()
	serverHandler := NewServerHandler()

	r := mux.NewRouter()

	r.HandleFunc("/search", serverHandler.SearchHandler).Methods("POST")
	r.HandleFunc("/fetch", serverHandler.NewFetchHandler).Methods("GET")
	r.HandleFunc("/fetch/{userid}/{page}", serverHandler.FetchHandler).Methods("GET")

	logger.Info("Starting server")
	if err := http.ListenAndServe(":3000", r); err != nil {
		logger.WithError(err).Fatal("Failed to start server, exiting")
	}
}
