// Code generated by MockGen. DO NOT EDIT.
// Source: API.go

// Package mock_hendlers is a generated GoMock package.
package mock_hendlers

import (
	context "context"
	http "net/http"
	reflect "reflect"
	sync "sync"

	gomock "github.com/golang/mock/gomock"
	websocket "github.com/gorilla/websocket"
	addition "github.com/whitewolf185/fs-go-hw/project/cmd/addition"
	add_Conn "github.com/whitewolf185/fs-go-hw/project/pkg/hendlers/add_Conn"
)

// MockConnectionService is a mock of ConnectionService interface.
type MockConnectionService struct {
	ctrl     *gomock.Controller
	recorder *MockConnectionServiceMockRecorder
}

// MockConnectionServiceMockRecorder is the mock recorder for MockConnectionService.
type MockConnectionServiceMockRecorder struct {
	mock *MockConnectionService
}

// NewMockConnectionService creates a new mock instance.
func NewMockConnectionService(ctrl *gomock.Controller) *MockConnectionService {
	mock := &MockConnectionService{ctrl: ctrl}
	mock.recorder = &MockConnectionServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConnectionService) EXPECT() *MockConnectionServiceMockRecorder {
	return m.recorder
}

// PingPong mocks base method.
func (m *MockConnectionService) PingPong(wg *sync.WaitGroup, ctx context.Context, ws *websocket.Conn) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PingPong", wg, ctx, ws)
}

// PingPong indicates an expected call of PingPong.
func (mr *MockConnectionServiceMockRecorder) PingPong(wg, ctx, ws interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PingPong", reflect.TypeOf((*MockConnectionService)(nil).PingPong), wg, ctx, ws)
}

// PrepareCandles mocks base method.
func (m *MockConnectionService) PrepareCandles(ws *websocket.Conn, wg *sync.WaitGroup, ctx context.Context, options addition.Options) (chan add_Conn.EventMsg, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareCandles", ws, wg, ctx, options)
	ret0, _ := ret[0].(chan add_Conn.EventMsg)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrepareCandles indicates an expected call of PrepareCandles.
func (mr *MockConnectionServiceMockRecorder) PrepareCandles(ws, wg, ctx, options interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareCandles", reflect.TypeOf((*MockConnectionService)(nil).PrepareCandles), ws, wg, ctx, options)
}

// Unsubscribe mocks base method.
func (m *MockConnectionService) Unsubscribe(ws *websocket.Conn) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unsubscribe", ws)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unsubscribe indicates an expected call of Unsubscribe.
func (mr *MockConnectionServiceMockRecorder) Unsubscribe(ws interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unsubscribe", reflect.TypeOf((*MockConnectionService)(nil).Unsubscribe), ws)
}

// WebConn mocks base method.
func (m *MockConnectionService) WebConn(url string) (*websocket.Conn, *http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WebConn", url)
	ret0, _ := ret[0].(*websocket.Conn)
	ret1, _ := ret[1].(*http.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// WebConn indicates an expected call of WebConn.
func (mr *MockConnectionServiceMockRecorder) WebConn(url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WebConn", reflect.TypeOf((*MockConnectionService)(nil).WebConn), url)
}
