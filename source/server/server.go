package server

import (
	"net/http"

	"github.com/ashmeet13/YoutubeDataService/source/common"
)

func Start() {
	logger := common.GetLogger()
	http.HandleFunc("/search", SearchHandler)

	logger.Info("Starting server")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.WithError(err).Fatal("Failed to start server, exiting")
	}

}
