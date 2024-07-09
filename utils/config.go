package utils

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DbUrl     string `yaml:"dbUrl"`
	DbUser    string `yaml:"dbUser"`
	DbPass    string `yaml:"dbPass"`
	Websocket struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}
}

func LoadConfig(configPath string) Config {
	var config Config

	yamlFile, err := ioutil.ReadFile(configPath)

	if err != nil {
		LogFatal("Error reading YAML file, %v\n", "red")
	}

	err = yaml.Unmarshal(yamlFile, &config)

	if err != nil {
		LogFatal("Error unmarshalling YAML: %v\n ", "red")
	}

	return config
}
