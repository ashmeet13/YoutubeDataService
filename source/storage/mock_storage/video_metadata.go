// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ashmeet13/YoutubeDataService/source/storage (interfaces: VideoMetadataInterface)

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	reflect "reflect"
	time "time"

	storage "github.com/ashmeet13/YoutubeDataService/source/storage"
	gomock "github.com/golang/mock/gomock"
)

// MockVideoMetadataInterface is a mock of VideoMetadataInterface interface.
type MockVideoMetadataInterface struct {
	ctrl     *gomock.Controller
	recorder *MockVideoMetadataInterfaceMockRecorder
}

// MockVideoMetadataInterfaceMockRecorder is the mock recorder for MockVideoMetadataInterface.
type MockVideoMetadataInterfaceMockRecorder struct {
	mock *MockVideoMetadataInterface
}

// NewMockVideoMetadataInterface creates a new mock instance.
func NewMockVideoMetadataInterface(ctrl *gomock.Controller) *MockVideoMetadataInterface {
	mock := &MockVideoMetadataInterface{ctrl: ctrl}
	mock.recorder = &MockVideoMetadataInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVideoMetadataInterface) EXPECT() *MockVideoMetadataInterfaceMockRecorder {
	return m.recorder
}

// BulkInsertMetadata mocks base method.
func (m *MockVideoMetadataInterface) BulkInsertMetadata(arg0 []*storage.VideoMetadata) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BulkInsertMetadata", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// BulkInsertMetadata indicates an expected call of BulkInsertMetadata.
func (mr *MockVideoMetadataInterfaceMockRecorder) BulkInsertMetadata(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BulkInsertMetadata", reflect.TypeOf((*MockVideoMetadataInterface)(nil).BulkInsertMetadata), arg0)
}

// FetchPagedMetadata mocks base method.
func (m *MockVideoMetadataInterface) FetchPagedMetadata(arg0 time.Time, arg1, arg2 int64) ([]*storage.VideoMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchPagedMetadata", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*storage.VideoMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchPagedMetadata indicates an expected call of FetchPagedMetadata.
func (mr *MockVideoMetadataInterfaceMockRecorder) FetchPagedMetadata(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchPagedMetadata", reflect.TypeOf((*MockVideoMetadataInterface)(nil).FetchPagedMetadata), arg0, arg1, arg2)
}

// FindOneMetadata mocks base method.
func (m *MockVideoMetadataInterface) FindOneMetadata(arg0, arg1 string) (*storage.VideoMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneMetadata", arg0, arg1)
	ret0, _ := ret[0].(*storage.VideoMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneMetadata indicates an expected call of FindOneMetadata.
func (mr *MockVideoMetadataInterfaceMockRecorder) FindOneMetadata(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneMetadata", reflect.TypeOf((*MockVideoMetadataInterface)(nil).FindOneMetadata), arg0, arg1)
}

// FindOneMetadataWithVideoID mocks base method.
func (m *MockVideoMetadataInterface) FindOneMetadataWithVideoID(arg0 string) (*storage.VideoMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneMetadataWithVideoID", arg0)
	ret0, _ := ret[0].(*storage.VideoMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneMetadataWithVideoID indicates an expected call of FindOneMetadataWithVideoID.
func (mr *MockVideoMetadataInterfaceMockRecorder) FindOneMetadataWithVideoID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneMetadataWithVideoID", reflect.TypeOf((*MockVideoMetadataInterface)(nil).FindOneMetadataWithVideoID), arg0)
}

// UpdateOneMetadata mocks base method.
func (m *MockVideoMetadataInterface) UpdateOneMetadata(arg0 string, arg1 *storage.VideoMetadata) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOneMetadata", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOneMetadata indicates an expected call of UpdateOneMetadata.
func (mr *MockVideoMetadataInterfaceMockRecorder) UpdateOneMetadata(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOneMetadata", reflect.TypeOf((*MockVideoMetadataInterface)(nil).UpdateOneMetadata), arg0, arg1)
}
