package server

import (
	"fmt"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/prometheus/client_golang/prometheus"
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
	customRegistry := prometheus.NewRegistry() // create custom registry for your custom metrics
	customCounter := prometheus.NewCounter(    // create new counter metric. This is replacement for `prometheus.Metric` struct
		prometheus.CounterOpts{
			Name: "custom_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
	)
	if err := customRegistry.Register(customCounter); err != nil { // register your new counter metric with metrics registry
		log.Fatal(err)
	}

	s.app.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		AfterNext: func(c echo.Context, err error) {
			customCounter.Inc() // use our custom metric in middleware. after every request increment the counter
		},
		Registerer: customRegistry, // use our custom registry instead of default Prometheus registry
	}))

	// Health check adding
	s.app.GET("/ping", func(c echo.Context) error {
		return c.String(200, "OK")
	})
	s.app.GET("/metrics", echoprometheus.NewHandlerWithConfig(echoprometheus.HandlerConfig{Gatherer: customRegistry}))
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
