package app

import (
	"github.com/ksopin/notes/pkg/db"
)

type Config struct{
	ProjectName string `yaml:"project_name"`
	Code string
	Version string
	Db *db.Config
	Auth *struct{
		KeyPath string `yaml:"key_path"`
	}
}

