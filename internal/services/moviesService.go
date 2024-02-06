package services

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"movieMatchMaker/models"
	"net/http"
	"net/url"
	"time"
)

const KpApiKey =
const baseApiURL = "https://api.kinopoisk.dev/v1.4/movie?page=1&limit=100&type=movie&rating.kp=7-10"

type Repository interface {
	AddMoviesToRoom(roomId string, movies []models.Movie)
	GetMoviesInRoom(roomId string) ([]models.Movie, error)
	RateMovie(roomId string, movieId int) ([]models.Movie, error)
}

type MoviesService struct {
	log        *slog.Logger
	repository Repository
}

func NewMoviesService(log *slog.Logger, r Repository) *MoviesService {
	return &MoviesService{
		log:        log,
		repository: r,
	}
}

func (s *MoviesService) GetMoviesInRoom(roomId string) ([]models.Movie, error) {
	movies, err := s.repository.GetMoviesInRoom(roomId)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (s *MoviesService) AddMoviesToRoom(roomId string, params models.Input) error {
	movies, err := getKpAPI(params)
	if err != nil {
		s.log.Error(fmt.Sprintf("kp api error: %s", err))
		return fmt.Errorf("kp api error: %s", err)
	}

	s.repository.AddMoviesToRoom(roomId, movies)
	s.log.Info(fmt.Sprintf("movies have been successfully added to the room <%s>", roomId))

	return nil
}

func (s *MoviesService) RateMovie(roomId string, movieId int) ([]models.Movie, error) {
	movie, err := s.repository.RateMovie(roomId, movieId)
	if err != nil {
		return []models.Movie{}, err
	}

	return movie, nil
}

func getKpAPI(params models.Input) ([]models.Movie, error) {
	URL := baseApiURL
	for _, years := range params.Years {
		URL += fmt.Sprintf("&year=%s", years)
	}
	for _, genre := range params.Genres {
		URL += fmt.Sprintf("&genres.name=%s", url.QueryEscape(genre))
	}

	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("X-API-KEY", KpApiKey)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	input := struct {
		Docs []models.Movie `json:"docs"`
	}{}
	err := json.NewDecoder(res.Body).Decode(&input)
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(input.Docs),
		func(i, j int) { input.Docs[i], input.Docs[j] = input.Docs[j], input.Docs[i] })

	return input.Docs[:30], nil
}
