package main

import (
	"Avito/internal/handler"
	"github.com/spf13/viper"
	"log"
)

// @title API-Server AvitoTech
// @version 1.0
// @description Microservice for working with user balance

// @host localhost:8080
func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("[Main] Error initialize configs: %s", err.Error())
	}

	srv := handler.NewServer()
	if err := srv.Start(viper.GetString("port")); err != nil {
		log.Fatalf("[Main] Server start error: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	return viper.ReadInConfig()
}
