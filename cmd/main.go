package main

import (
	"Avito/internal/handler"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
)

// @title Billing App API
// @version 1.0
// @description API server for Avito billing app

// @host localhost:8080
// @BasePath /

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("Can't initialize configs: %s", err.Error())
	}

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Can't load env variables: %s", err.Error())
	}

	srv := handler.NewServer()
	if err := srv.Start(viper.GetString("port")); err != nil {
		log.Fatalf("Can't get the server running: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	return viper.ReadInConfig()
}
