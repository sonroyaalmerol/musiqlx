package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigFile(".env")
	/*if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}*/
	log.Println("Config loaded successfully.")
}
