package config

import (
	"github.com/ksopin/notes/internal/app"
	"github.com/ksopin/notes/internal/http"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	App *app.Config
	Http *http.Config
}


func GetConfig(configFile string) *Config {
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("yamlFile.Get %+v ", err)
	}
	c := &Config{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %+v", err)
	}

	return c
}
