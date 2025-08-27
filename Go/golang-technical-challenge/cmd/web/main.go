package main

import (
	"fmt"
	"golang-technical-challenge/internal/config"
)

func main() {
	v := config.NewViper()
	log := config.NewLogger(v)
	db := config.NewDatabase(v, log)
	app := config.NewFiber(v)
	validate := config.NewValidator(v)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   v,
	})

	webPort := v.GetInt("APP_PORT")
	log.Infof("Starting server on port %d ...", webPort)

	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
