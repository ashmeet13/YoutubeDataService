package common

import (
	"os"
	"strconv"
	"strings"
)

const (
	MongoBaseURL      = "MONGO_BASE_URL"
	MongoDatabaseName = "MONGO_DATABASE_NAME"
	YoutubeAPIKeys    = "YOUTUBE_API_KEYS"
	DefaultPageSize   = "DEFAULT_PAGE_SIZE"
)

type Configuration struct {
	MongoBaseURL      string
	MongoDatabaseName string
	YoutubeAPIKeys    []string
	DefaultPageSize   int
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

	defaultPageSizeString := os.Getenv(DefaultPageSize)
	if defaultPageSizeString == "" {
		defaultPageSizeString = "20"
	}

	keys := []string{}
	for _, apiKey := range strings.Split(youtubeAPIKeys, ",") {
		keys = append(keys, apiKey)
	}

	defaultPageSize, err := strconv.Atoi(defaultPageSizeString)
	if err != nil {
		return nil
	}

	return &Configuration{
		MongoBaseURL:      mongoBaseURL,
		MongoDatabaseName: mongoDatabaseName,
		YoutubeAPIKeys:    keys,
		DefaultPageSize:   defaultPageSize,
	}
}
