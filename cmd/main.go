package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"movieMatchMaker/internal/config"
	"movieMatchMaker/internal/handlers"
	"movieMatchMaker/internal/repository"
	"movieMatchMaker/internal/services"
	"movieMatchMaker/pkg/logger"
	"movieMatchMaker/pkg/middleware"
	"movieMatchMaker/pkg/server"
	"os"
	"os/signal"
	"syscall"
)

// @title movieMatchMaker API
// @version 1.0
// @description API Server

// @host localhost:8080
// @BasePath /

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)
	log.Debug("init logger successful", slog.String("env", cfg.Env))

	srv := server.Init()
	log.Debug("init server successful")

	api := InitRouter(cfg, log)
	log.Debug("init router successful")

	roomsRepo := repository.NewRoomsRepository()
	roomsService := services.NewRoomService(log, roomsRepo)
	roomsHandler := handlers.NewRoomsHandler(roomsService)
	roomsHandler.Register(api)

	moviesService := services.NewMoviesService(log, roomsRepo)
	moviesHandler := handlers.NewMovieHandler(moviesService)
	moviesHandler.Register(api)

	go roomsRepo.CleanupRooms()

	log.Debug("init app successful")
	go func() {
		err := srv.Run(cfg, api)
		if err != nil {
			log.Error("failed to run server")
		}
	}()
	log.Info("app started", slog.String("address", cfg.Address))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("app shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Error(fmt.Sprintf("error occured on server shutting down: %s", err.Error()))
	}
}

func InitRouter(cfg config.Config, log *slog.Logger) *gin.Engine {
	gin.SetMode(cfg.Env)
	router := gin.New()
	router.Use(middleware.Logging(log))

	return router
}
