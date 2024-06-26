package services

import (
	"github.com/muety/wakapi/config"
	"github.com/muety/wakapi/repositories"
	"time"

	"github.com/muety/wakapi/models"
)

type HeartbeatService struct {
	config              *config.Config
	repository          repositories.IHeartbeatRepository
	languageMappingSrvc ILanguageMappingService
}

func NewHeartbeatService(heartbeatRepo repositories.IHeartbeatRepository, languageMappingService ILanguageMappingService) *HeartbeatService {
	return &HeartbeatService{
		config:              config.Get(),
		repository:          heartbeatRepo,
		languageMappingSrvc: languageMappingService,
	}
}

func (srv *HeartbeatService) InsertBatch(heartbeats []*models.Heartbeat) error {
	return srv.repository.InsertBatch(heartbeats)
}

func (srv *HeartbeatService) GetAllWithin(from, to time.Time, user *models.User) ([]*models.Heartbeat, error) {
	heartbeats, err := srv.repository.GetAllWithin(from, to, user)
	if err != nil {
		return nil, err
	}
	return srv.augmented(heartbeats, user.ID)
}

func (srv *HeartbeatService) GetFirstByUsers() ([]*models.TimeByUser, error) {
	return srv.repository.GetFirstByUsers()
}

func (srv *HeartbeatService) DeleteBefore(t time.Time) error {
	return srv.repository.DeleteBefore(t)
}

func (srv *HeartbeatService) augmented(heartbeats []*models.Heartbeat, userId string) ([]*models.Heartbeat, error) {
	languageMapping, err := srv.languageMappingSrvc.ResolveByUser(userId)
	if err != nil {
		return nil, err
	}

	for i := range heartbeats {
		heartbeats[i].Augment(languageMapping)
	}

	return heartbeats, nil
}
