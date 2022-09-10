package main

import "github.com/ashmeet13/YoutubeDataService/common"

func main() {
	config := common.GetConfiguration()
	logger := common.GetLogger()

	logger.
		WithField("MongoURL", config.MongoBaseURL).
		WithField("MongoDatabase", config.MongoDatabaseName).
		Info("Initalising Server")
}
