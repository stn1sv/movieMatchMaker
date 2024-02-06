package repository

import (
	"fmt"
	"movieMatchMaker/models"
	"sync"
	"time"
)

type RoomsRepository struct {
	rooms map[string]*models.Room
	m     sync.RWMutex
}

func NewRoomsRepository() *RoomsRepository {
	return &RoomsRepository{
		rooms: make(map[string]*models.Room),
	}
}

func (r *RoomsRepository) CreateRoom(userName string, roomId string) []models.User {
	user := models.User{
		UserName: userName,
		IsAdmin:  true,
	}
	room := &models.Room{
		Id:           roomId,
		Participants: []models.User{user},
		CreatedAt:    time.Now().Unix(),
	}

	r.m.Lock()
	defer r.m.Unlock()
	r.rooms[roomId] = room

	return room.Participants
}

func (r *RoomsRepository) JoinToRoom(userName string, roomId string) ([]models.User, error) {
	r.m.Lock()
	defer r.m.Unlock()

	room, ok := r.rooms[roomId]
	if !ok {
		return nil, fmt.Errorf("room not found")
	}

	room.Participants = append(room.Participants, models.User{UserName: userName, IsAdmin: false})

	return room.Participants, nil
}

func (r *RoomsRepository) GetAllParticipants(roomId string) ([]models.User, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	room, ok := r.rooms[roomId]
	if !ok {
		return nil, fmt.Errorf("room not found")
	}

	return room.Participants, nil
}

func (r *RoomsRepository) CleanupRooms() {
	for {
		time.Sleep(time.Minute * 30)
		currentTimestamp := time.Now().Unix()

		r.m.Lock()
		for roomID, room := range r.rooms {
			if currentTimestamp-room.CreatedAt > 1800 {
				delete(r.rooms, roomID)
			}
		}
		r.m.Unlock()
	}
}

func (r *RoomsRepository) GetMoviesInRoom(roomId string) ([]models.Movie, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	room, ok := r.rooms[roomId]
	if !ok {
		return nil, fmt.Errorf("room not found")
	}

	return room.Movies, nil
}

func (r *RoomsRepository) AddMoviesToRoom(roomId string, movies []models.Movie) {
	r.m.Lock()
	defer r.m.Unlock()
	_, ok := r.rooms[roomId]
	if !ok {
		return
	}

	r.rooms[roomId].Movies = movies
}

func (r *RoomsRepository) RateMovie(roomId string, movieId int) ([]models.Movie, error) {
	r.m.Lock()
	defer r.m.Unlock()

	room, ok := r.rooms[roomId]
	if !ok {
		return []models.Movie{}, fmt.Errorf("room not found")
	}

	room.MoviesThatEveryoneLiked = []models.Movie{}
	for i, _ := range room.Movies {
		if room.Movies[i].Id == movieId {
			room.Movies[i].UsersWhoLikedItCount++
		}
		if room.Movies[i].UsersWhoLikedItCount == len(room.Participants) {
			room.MoviesThatEveryoneLiked = append(room.MoviesThatEveryoneLiked, room.Movies[i])
		}
	}
	r.rooms[roomId] = room

	return room.MoviesThatEveryoneLiked, nil
}
