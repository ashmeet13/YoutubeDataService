package common

import "os"

const (
	MongoBaseURL      = "MONGO_BASE_URL"
	MongoDatabaseName = "MONGO_DATABASE_NAME"
	YoutubeAPIKey     = "YOUTUBE_API_KEY"
)

type Configuration struct {
	MongoBaseURL      string
	MongoDatabaseName string
	YoutubeAPIKey     string
}

var Config *Configuration

func GetConfiguration() *Configuration {
	if Config == nil {
		Config = SetupConfiguration()
	}
	return Config
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

	youtubeAPIKey := os.Getenv(YoutubeAPIKey)
	if youtubeAPIKey == "" {
		logger.Fatalln("Could not find environment variable", YoutubeAPIKey)
		return nil
	}

	return &Configuration{
		MongoBaseURL:      mongoBaseURL,
		MongoDatabaseName: mongoDatabaseName,
		YoutubeAPIKey:     youtubeAPIKey,
	}
}
