// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/rafaelrubbioli/fileapi/pkg/entity"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockService) Create(ctx context.Context, user, size int, name, path, contentType string, file io.Reader) (*entity.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user, size, name, path, contentType, file)
	ret0, _ := ret[0].(*entity.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockServiceMockRecorder) Create(ctx, user, size, name, path, contentType, file interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockService)(nil).Create), ctx, user, size, name, path, contentType, file)
}

// Delete mocks base method.
func (m *MockService) Delete(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockServiceMockRecorder) Delete(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockService)(nil).Delete), ctx, key)
}

// Get mocks base method.
func (m *MockService) Get(ctx context.Context, id string) (*entity.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*entity.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockServiceMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockService)(nil).Get), ctx, id)
}

// GetByUser mocks base method.
func (m *MockService) GetByUser(ctx context.Context, user int, prefix string) ([]*entity.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUser", ctx, user, prefix)
	ret0, _ := ret[0].([]*entity.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUser indicates an expected call of GetByUser.
func (mr *MockServiceMockRecorder) GetByUser(ctx, user, prefix interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUser", reflect.TypeOf((*MockService)(nil).GetByUser), ctx, user, prefix)
}

// Move mocks base method.
func (m *MockService) Move(ctx context.Context, user int, id, newPath string) (*entity.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Move", ctx, user, id, newPath)
	ret0, _ := ret[0].(*entity.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Move indicates an expected call of Move.
func (mr *MockServiceMockRecorder) Move(ctx, user, id, newPath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Move", reflect.TypeOf((*MockService)(nil).Move), ctx, user, id, newPath)
}
