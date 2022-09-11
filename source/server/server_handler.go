package server

import (
	"sync"

	"github.com/ashmeet13/YoutubeDataService/source/storage"
)

func NewServerHandler() *ServerHandler {
	return &ServerHandler{
		lock:              sync.Mutex{},
		userLastSentIndex: map[string]int{},
		userPageSize:      map[string]int{},
		storageHandler:    storage.NewVideoMetadataImpl(),
	}
}

type ServerHandler struct {
	lock              sync.Mutex
	userLastSentIndex map[string]int
	userPageSize      map[string]int
	storageHandler    storage.VideoMetadataInterface
}
