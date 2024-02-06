package services

import (
	"fmt"
	"log/slog"
	"movieMatchMaker/models"
	"movieMatchMaker/pkg/random"
)

type RoomsRepository interface {
	CreateRoom(userName string, roomId string) []models.User
	JoinToRoom(userName string, roomId string) ([]models.User, error)
	GetAllParticipants(roomId string) ([]models.User, error)
}

type RoomsService struct {
	log        *slog.Logger
	repository RoomsRepository
}

func NewRoomService(log *slog.Logger, r RoomsRepository) *RoomsService {
	return &RoomsService{
		log:        log,
		repository: r,
	}
}

func (s *RoomsService) CreateRoom(userName string) (string, []models.User) {
	roomId := random.GenerateRandomID()

	participants := s.repository.CreateRoom(userName, roomId)
	s.log.Info(fmt.Sprintf("room <%s> has been created by user <%s>", roomId, userName))

	return roomId, participants
}

func (s *RoomsService) JoinToRoom(userName string, roomId string) ([]models.User, error) {
	participants, err := s.repository.JoinToRoom(userName, roomId)
	if err != nil {
		return nil, err
	}
	s.log.Info(fmt.Sprintf("user <%s> join to room <%s>", userName, roomId))

	return participants, nil
}

func (s *RoomsService) GetAllParticipant(roomId string) ([]models.User, error) {
	participants, err := s.repository.GetAllParticipants(roomId)
	if err != nil {
		return nil, err
	}

	return participants, nil
}
