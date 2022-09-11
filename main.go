package main

import (
	"github.com/ashmeet13/YoutubeDataService/source/common"
	"github.com/ashmeet13/YoutubeDataService/source/server"
	"github.com/ashmeet13/YoutubeDataService/source/worker"
)

func main() {
	config := common.GetConfiguration()
	logger := common.GetLogger()

	logger.
		WithField("MongoURL", config.MongoBaseURL).
		WithField("MongoDatabase", config.MongoDatabaseName).
		Info("Initalising Server")

	workerHandler, err := worker.NewWorkerHandler(config.YoutubeAPIKeys)

	if err != nil {
		logger.Fatal("Failed to init worker")
	}

	go workerHandler.Start()
	server.Start()
}
