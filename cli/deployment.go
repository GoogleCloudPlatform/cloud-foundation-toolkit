package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Name        string
	Project     string
	file_path   string
	yaml_string string
}

func NewConfig(file_path string) *Config {
	data, err := ioutil.ReadFile(file_path)
	if err != nil {
		log.Fatal(err)
	}
	yaml_string := string(data)

	config := &Config{
		file_path:   file_path,
		yaml_string: yaml_string,
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return config
}
