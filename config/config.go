package config

import (
	"encoding/json"
	"io/ioutil"
)

type MongoConfig struct {
	Addr string `json:"addr"`
	Port uint   `json:"port"`
	DB   string `json:"db"`
}
type MQTTConfig struct {
	Addr string `json:"addr"`
	Port uint   `json:"port"`
	Name string `json:"name"`
}

type Config struct {
	Mongo MongoConfig `json:"mongo"`
	MQTT  MQTTConfig  `json:"mqtt"`
}

func LoadConfig(config *Config) error {
	data, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		return err
	}
	json.Unmarshal(data, config)
	if err != nil {
		return err
	}
	return nil

}
