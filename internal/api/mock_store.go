// Code generated by MockGen. DO NOT EDIT.
// Source: internal/api/api.go

// Package mock_api is a generated GoMock package.
package api

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	go_uuid "github.com/satori/go.uuid"
)

// MockStorager is a mock of Storager interface.
type MockStorager struct {
	ctrl     *gomock.Controller
	recorder *MockStoragerMockRecorder
}

// MockStoragerMockRecorder is the mock recorder for MockStorager.
type MockStoragerMockRecorder struct {
	mock *MockStorager
}

// NewMockStorager creates a new mock instance.
func NewMockStorager(ctrl *gomock.Controller) *MockStorager {
	mock := &MockStorager{ctrl: ctrl}
	mock.recorder = &MockStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorager) EXPECT() *MockStoragerMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStorager) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockStoragerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorager)(nil).Close))
}

// DeleteRecords mocks base method.
func (m *MockStorager) DeleteRecords(ctx context.Context, deleteItems []DeleteItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRecords", ctx, deleteItems)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRecords indicates an expected call of DeleteRecords.
func (mr *MockStoragerMockRecorder) DeleteRecords(ctx, deleteItems interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRecords", reflect.TypeOf((*MockStorager)(nil).DeleteRecords), ctx, deleteItems)
}

// Insert mocks base method.
func (m *MockStorager) Insert(ctx context.Context, userID go_uuid.UUID, key, value string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, userID, key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockStoragerMockRecorder) Insert(ctx, userID, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockStorager)(nil).Insert), ctx, userID, key, value)
}

// InsertBatch mocks base method.
func (m *MockStorager) InsertBatch(arg0 context.Context, userID go_uuid.UUID, batch []BatchElement) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertBatch", arg0, userID, batch)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertBatch indicates an expected call of InsertBatch.
func (mr *MockStoragerMockRecorder) InsertBatch(arg0, userID, batch interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertBatch", reflect.TypeOf((*MockStorager)(nil).InsertBatch), arg0, userID, batch)
}

// Ping mocks base method.
func (m *MockStorager) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockStoragerMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockStorager)(nil).Ping), ctx)
}

// Select mocks base method.
func (m *MockStorager) Select(ctx context.Context, key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Select", ctx, key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Select indicates an expected call of Select.
func (mr *MockStoragerMockRecorder) Select(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Select", reflect.TypeOf((*MockStorager)(nil).Select), ctx, key)
}

// SelectUserAll mocks base method.
func (m *MockStorager) SelectUserAll(ctx context.Context, userID go_uuid.UUID) ([]BatchElement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectUserAll", ctx, userID)
	ret0, _ := ret[0].([]BatchElement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectUserAll indicates an expected call of SelectUserAll.
func (mr *MockStoragerMockRecorder) SelectUserAll(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectUserAll", reflect.TypeOf((*MockStorager)(nil).SelectUserAll), ctx, userID)
}
