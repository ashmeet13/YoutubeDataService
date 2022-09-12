package worker

import (
	"testing"

	"github.com/ashmeet13/YoutubeDataService/source/storage/mock_storage"
	"github.com/ashmeet13/YoutubeDataService/source/youtube/mock_youtube"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type WorkerHandlerSuite struct {
	suite.Suite
	*require.Assertions
	ctrl *gomock.Controller

	mockVideoMetadataStore *mock_storage.MockVideoMetadataInterface
	mockYoutubeHandler     *mock_youtube.MockYoutubeInterface

	workerHandler *WorkerHandler
}

func TestWorkerHandlerSuite(t *testing.T) {
	suite.Run(t, new(WorkerHandlerSuite))
}

func (s *WorkerHandlerSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.ctrl = gomock.NewController(s.T())

	s.mockVideoMetadataStore = mock_storage.NewMockVideoMetadataInterface(s.ctrl)
	s.mockYoutubeHandler = mock_youtube.NewMockYoutubeInterface(s.ctrl)

	s.workerHandler = &WorkerHandler{
		videoMetadataHandler: s.mockVideoMetadataStore,
		youtubeHandler:       s.mockYoutubeHandler,
	}
}

func (s *WorkerHandlerSuite) TearDownSuite() {
	s.ctrl.Finish()
}
