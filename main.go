package main

import (
	"fmt"

	"github.com/ashmeet13/YoutubeDataService/source/common"
	"github.com/ashmeet13/YoutubeDataService/source/worker"
)

func main() {
	config := common.GetConfiguration()
	logger := common.GetLogger()

	logger.
		WithField("MongoURL", config.MongoBaseURL).
		WithField("MongoDatabase", config.MongoDatabaseName).
		Info("Initalising Server")

	workerHandler := worker.NewWorkerHandler(config.YoutubeAPIKeys)
	fmt.Println(workerHandler)

	workerHandler.Start()
}
