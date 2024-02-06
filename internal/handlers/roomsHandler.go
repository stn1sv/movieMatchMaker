package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"movieMatchMaker/models"
	"net/http"
)

type RoomService interface {
	CreateRoom(userName string) (string, []models.User)
	JoinToRoom(userName string, roomId string) ([]models.User, error)
	GetAllParticipant(roomId string) ([]models.User, error)
}

type RoomsHandler struct {
	service RoomService
}

func NewRoomsHandler(s RoomService) *RoomsHandler {
	return &RoomsHandler{
		service: s,
	}
}

func (h *RoomsHandler) Register(api *gin.Engine) {
	api.GET("/getAllParticipants", h.GetAllParticipants)
	api.GET("/createRoom", h.CreateRoom)
	api.GET("/joinRoom", h.JoinToRoom)
}

func (h *RoomsHandler) CreateRoom(c *gin.Context) {
	userName := c.Query("userName")
	if userName == "undefined" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("empty userName"))
		return
	}

	roomId, participants := h.service.CreateRoom(userName)

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{"roomId": roomId, "participants": participants})
}

func (h *RoomsHandler) JoinToRoom(c *gin.Context) {
	userName := c.Query("userName")
	roomId := c.Query("roomId")
	if userName == "undefined" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("empty userName"))
		return
	}
	if roomId == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("empty roomId"))
		return
	}

	participants, err := h.service.JoinToRoom(userName, roomId)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{"participants": participants})
}

func (h *RoomsHandler) GetAllParticipants(c *gin.Context) {
	roomId := c.Query("roomId")
	if roomId == "" {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("empty roomId"))
		return
	}

	participants, err := h.service.GetAllParticipant(roomId)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{"participants": participants})
}
