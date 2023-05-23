// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockNotificationApp is a mock of NotificationApp interface.
type MockNotificationApp struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationAppMockRecorder
}

// MockNotificationAppMockRecorder is the mock recorder for MockNotificationApp.
type MockNotificationAppMockRecorder struct {
	mock *MockNotificationApp
}

// NewMockNotificationApp creates a new mock instance.
func NewMockNotificationApp(ctrl *gomock.Controller) *MockNotificationApp {
	mock := &MockNotificationApp{ctrl: ctrl}
	mock.recorder = &MockNotificationAppMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotificationApp) EXPECT() *MockNotificationAppMockRecorder {
	return m.recorder
}

// AddUserToNotificationTopic mocks base method.
func (m *MockNotificationApp) AddUserToNotificationTopic(topic string, token models.NotificationToken, ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserToNotificationTopic", topic, token, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserToNotificationTopic indicates an expected call of AddUserToNotificationTopic.
func (mr *MockNotificationAppMockRecorder) AddUserToNotificationTopic(topic, token, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserToNotificationTopic", reflect.TypeOf((*MockNotificationApp)(nil).AddUserToNotificationTopic), topic, token, ctx)
}

// RemoveUserFromNotificationTopic mocks base method.
func (m *MockNotificationApp) RemoveUserFromNotificationTopic(topic string, token models.NotificationToken, ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUserFromNotificationTopic", topic, token, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUserFromNotificationTopic indicates an expected call of RemoveUserFromNotificationTopic.
func (mr *MockNotificationAppMockRecorder) RemoveUserFromNotificationTopic(topic, token, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUserFromNotificationTopic", reflect.TypeOf((*MockNotificationApp)(nil).RemoveUserFromNotificationTopic), topic, token, ctx)
}

// SendUserNotification mocks base method.
func (m *MockNotificationApp) SendUserNotification(notification models.Notification, ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendUserNotification", notification, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendUserNotification indicates an expected call of SendUserNotification.
func (mr *MockNotificationAppMockRecorder) SendUserNotification(notification, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendUserNotification", reflect.TypeOf((*MockNotificationApp)(nil).SendUserNotification), notification, ctx)
}
