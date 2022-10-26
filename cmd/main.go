package main

import (
	"Avito/internal/handler"
	"github.com/spf13/viper"
	"log"
)

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
