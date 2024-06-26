package mocks

import (
	"github.com/muety/wakapi/models"
	"github.com/stretchr/testify/mock"
	"time"
)

type HeartbeatServiceMock struct {
	mock.Mock
}

func (m *HeartbeatServiceMock) InsertBatch(heartbeats []*models.Heartbeat) error {
	args := m.Called(heartbeats)
	return args.Error(0)
}

func (m *HeartbeatServiceMock) GetAllWithin(time time.Time, time2 time.Time, user *models.User) ([]*models.Heartbeat, error) {
	args := m.Called(time, time2, user)
	return args.Get(0).([]*models.Heartbeat), args.Error(1)
}

func (m *HeartbeatServiceMock) GetFirstByUsers() ([]*models.TimeByUser, error) {
	args := m.Called()
	return args.Get(0).([]*models.TimeByUser), args.Error(1)
}

func (m *HeartbeatServiceMock) DeleteBefore(time time.Time) error {
	args := m.Called(time)
	return args.Error(0)
}
