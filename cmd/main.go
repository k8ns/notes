package main

import (
	"github.com/ksopin/notes/pkg/app"
	"github.com/ksopin/notes/pkg/db"
	"github.com/ksopin/notes/pkg/http"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	App *app.Config
	Db *db.Config
	Http *http.Config
}

func main() {

	cfg, err := ParseConfig("config.yml")
	if err != nil {
		panic(err)
	}

	db.InitConnection(cfg.Db)

	err = http.Run(cfg.Http, cfg.App)
	if err != nil {
		panic(err)
	}
}

func ParseConfig(configFile string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
