package server

import (
	"github.com/ashmeet13/YoutubeDataService/source/storage"
)

func NewServerHandler() *ServerHandler {
	return &ServerHandler{
		videoMetadataHandler: storage.NewVideoMetadataImpl(),
		userHandler:          storage.NewUserImpl(),
	}
}

type ServerHandler struct {
	videoMetadataHandler storage.VideoMetadataInterface
	userHandler          storage.UserInterface
}
