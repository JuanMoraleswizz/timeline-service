package main

import (
	"fmt"

	"uala.com/timeline-service/cmd/server"
	"uala.com/timeline-service/config"
	"uala.com/timeline-service/internal/database"
)

func main() {
	// Start the server
	conf := config.GetConfig()
	db := database.NewMariadbDatabase(conf)
	redis := database.NewRedisDatabase(conf)

	fmt.Println("Starting server on port", conf.Server.Port)
	server.NewServer(conf, db, redis).Start()
}
