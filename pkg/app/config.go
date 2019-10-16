package app

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Dbserver string
	Dbusername string
	Dbpassword string
	Dbname string
	Httpport string
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

func (c *Config) DbServer() string {
	return c.Dbserver
}

func (c *Config) DbUsername() string {
	return c.Dbusername
}

func (c *Config) DbName() string {
	return c.Dbname
}

func (c *Config) DbPassword() string {
	return c.Dbpassword
}

