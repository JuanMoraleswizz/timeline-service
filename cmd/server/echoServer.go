package server

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"uala.com/timeline-service/config"
	"uala.com/timeline-service/internal/database"
	"uala.com/timeline-service/internal/handler"
	"uala.com/timeline-service/internal/repository"
	"uala.com/timeline-service/internal/usescase"
)

type server struct {
	app   *echo.Echo
	db    database.Database
	conf  *config.Config
	redis database.RedisDatabase
}

func NewServer(conf *config.Config, db database.Database, redis database.RedisDatabase) *server {
	echoApp := echo.New()
	echoApp.Logger.SetLevel(log.DEBUG)
	return &server{
		app:   echoApp,
		db:    db,
		conf:  conf,
		redis: redis,
	}
}

func (s *server) Start() {
	s.app.Use(middleware.Recover())
	s.app.Use(middleware.Logger())

	// Health check adding
	s.app.GET("/ping", func(c echo.Context) error {
		return c.String(200, "OK")
	})
	s.initializeCockroachHttpHandler()

	serverUrl := fmt.Sprintf(":%d", s.conf.Server.Port)
	s.app.Logger.Fatal(s.app.Start(serverUrl))
}

func (s *server) initializeCockroachHttpHandler() {
	// Initialize all repositories
	timeLineRepository := repository.NewTimelineMariadbRepository(s.db, s.redis)

	// Initialize all usecases
	timeLineUsesCase := usescase.NewGetTimeLineImpl(timeLineRepository)

	// Initialize all handlers
	userHandler := handler.NewTimelineHandler(timeLineUsesCase)

	// Routers
	s.app.GET("/v2/timeline/:userId/:forceSync", userHandler.GetTimeline)

}
