package main

import (
	"encoding/json"
	"errors"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"regexp"
)

type Config struct {
	Server string `yaml:"server" json:"server"`
	Path   string `yaml:"path" json:"path"`
	Days   int    `yaml:"days" json:"days"`
}

var (
	showConfig = kingpin.Flag("config.show", "Show actual configuration and ends").Default("false").Bool()
	configFile = kingpin.Flag("config.file", "Configuration file default is \"cfg.yml\".").PlaceHolder("cfg.yml").Default("cfg.yml").String()
	path       = kingpin.Flag("path", "Path where store export json data").PlaceHolder("path").Default("./dump").String()
	server     = kingpin.Flag("server", "Prometheus server FQDN or IP address").PlaceHolder("server").Default("").String()
	config     = &Config{
		Server: "",
		Path:   "./dump",
		Days:   1,
	}
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func dirAccessible(path string) bool {

}

func (c *Config) LoadFile(filename string) error {
	if fileExists(filename) {
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		err = yaml.UnmarshalStrict(content, c)
		if err != nil {
			err = json.Unmarshal(content, c)
			if err != nil {
				return err
			}
		}
	}
	if len(*server) > 0 {
		c.Server = *server
	}
	if len(*path) > 0 {
		c.Path = *path
	}

	match, err := regexp.MatchString("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$", c.Server)
	if !match || err != nil {
		match, err = regexp.MatchString("^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$", c.Server)
		if !match || err != nil {
			return errors.New("defined Prometheus server address isn't valid FQDN or IP address")
		}
	}
	if len(c.Path) < 1 {
		c.Path = "./dump"
	}
	if c.Days < 1 || c.Days > 60 {
		return errors.New("defined days back not valid (1 - 60)")
	}
	if !dirExists(c.Path) {
		return errors.New("path not exists")
	}
	if !dirAccessible(c.Path) {
		return errors.New("path not accessible for write")
	}
	return nil
}
