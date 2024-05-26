package main

import (
	"encoding/json"
	"os"

	"github.com/uncle-gua/log"
)

var config Config

type Config struct {
	Amount    float64 `json:"amount"`
	ApiKey    string  `json:"apiKey"`
	ApiSecret string  `json:"apiSecret"`
	Duration  int     `json:"duration"`
}

func init() {
	b, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(b, &config); err != nil {
		log.Fatal(err)
	}
	log.Info(config)
}
