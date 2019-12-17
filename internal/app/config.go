package app

import (
	"github.com/ksopin/notes/pkg/db"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)
type Config struct {
	App *AppConfig
	Db *db.Config
	Http *struct {
		Enabled bool
		Port int
	}
}

type AppConfig struct{
	ProjectName string `yaml:"project_name"`
	Code string
	Version string
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
