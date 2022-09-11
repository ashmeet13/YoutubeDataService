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

	r.HandleFunc("/search", serverHandler.SearchHandler)
	r.HandleFunc("/fetch", serverHandler.NewFetchHandler)
	r.HandleFunc("/fetch/{userid}/{page}", serverHandler.FetchHandler)

	logger.Info("Starting server")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.WithError(err).Fatal("Failed to start server, exiting")
	}
}