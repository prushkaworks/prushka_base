package server

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	ConnectionString string
	ServerPort       string
}

func configInit() *Config {
	conf := &Config{}
	rootDir, _ := os.Getwd()
	jsonFile, _ := ioutil.ReadFile(rootDir + "/config.template.json")
	json.Unmarshal(jsonFile, &conf)
	return conf
}
