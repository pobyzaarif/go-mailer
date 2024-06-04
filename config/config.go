package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Provider          string `json:"provider"`
	Host              string `json:"host"`
	Port              int    `json:"port"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	TestMailTo        string `json:"test_mail_to"`
	TestMailReaderURL string `json:"test_mail_reader_url"`
}

func LoadConfig(filename string) *Config {
	// Read the JSON file
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to read config file: %v", err))
	}

	// Unmarshal the JSON data into the config struct
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to unmarshal config JSON: %v", err))
	}

	return &config
}
