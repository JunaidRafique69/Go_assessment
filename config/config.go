package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Url    string `yaml:"url"`
	DbName string `yaml:"dbname"`
}

func LoadConfig() (DatabaseConfig, error) {

	viper.AddConfigPath("./config")
	viper.SetConfigName("database-config")
	viper.SetConfigType("yaml")
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	// Read the configurations file
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	var dbConfig DatabaseConfig

	// Unmarshal the config into the DatabaseConfig struct
	err = viper.Unmarshal(&dbConfig)
	if err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	return dbConfig, err
}
