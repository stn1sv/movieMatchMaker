package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"movieMatchMaker/models"
	"net/http"
	"strconv"
)

type MoviesService interface {
	GetMoviesInRoom(roomId string) ([]models.Movie, error)
	AddMoviesToRoom(roomId string, params models.Input) error
	RateMovie(roomId string, movieId int) ([]models.Movie, error)
}

type MovieHandler struct {
	service MoviesService
}

func NewMovieHandler(s MoviesService) *MovieHandler {
	return &MovieHandler{
		service: s,
	}
}

func (h *MovieHandler) Register(api *gin.Engine) {
	api.GET("/getMovies", h.GetMovies)
	api.POST("/setMovies", h.SetMovies)
	api.GET("/rateMovie", h.RateMovie)
}

func (h *MovieHandler) GetMovies(c *gin.Context) {
	roomId := c.Query("roomId")
	if roomId == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("empty roomId"))
		return
	}

	response, err := h.service.GetMoviesInRoom(roomId)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if response == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("empty movies"))
		return
	}
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, response)
}

func (h *MovieHandler) SetMovies(c *gin.Context) {
	roomId := c.Query("roomId")
	if roomId == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("empty roomId"))
		return
	}

	var input models.Input
	err := c.BindJSON(&input)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = h.service.AddMoviesToRoom(roomId, input)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.Status(http.StatusOK)
}

func (h *MovieHandler) RateMovie(c *gin.Context) {
	roomId := c.Query("roomId")
	if roomId == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("empty roomId"))
		return
	}
	movieId := c.Query("movieId")
	intMovieId, err := strconv.Atoi(movieId)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("empty roomId"))
		return
	}

	response, err := h.service.RateMovie(roomId, intMovieId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, response)
}
