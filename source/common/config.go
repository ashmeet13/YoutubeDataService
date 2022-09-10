package common

import (
	"os"
	"strings"
)

const (
	MongoBaseURL      = "MONGO_BASE_URL"
	MongoDatabaseName = "MONGO_DATABASE_NAME"
	YoutubeAPIKeys    = "YOUTUBE_API_KEYS"
)

type Configuration struct {
	MongoBaseURL      string
	MongoDatabaseName string
	YoutubeAPIKeys    []string
}

var config *Configuration

func GetConfiguration() *Configuration {
	if config == nil {
		config = SetupConfiguration()
	}
	return config
}

func SetupConfiguration() *Configuration {
	logger := GetLogger()

	mongoBaseURL := os.Getenv(MongoBaseURL)
	if mongoBaseURL == "" {
		logger.Fatalln("Could not find environment variable", MongoBaseURL)
		return nil
	}

	mongoDatabaseName := os.Getenv(MongoDatabaseName)
	if mongoDatabaseName == "" {
		logger.Fatalln("Could not find environment variable", MongoDatabaseName)
		return nil
	}

	youtubeAPIKeys := os.Getenv(YoutubeAPIKeys)
	if youtubeAPIKeys == "" {
		logger.Fatalln("Could not find environment variable", YoutubeAPIKeys)
		return nil
	}

	keys := []string{}
	for _, apiKey := range strings.Split(youtubeAPIKeys, ",") {
		keys = append(keys, apiKey)
	}

	return &Configuration{
		MongoBaseURL:      mongoBaseURL,
		MongoDatabaseName: mongoDatabaseName,
		YoutubeAPIKeys:    keys,
	}
}
