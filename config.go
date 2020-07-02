package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const connectionTimeout = 10

type configuration struct {
	Server string   `yaml:"server" json:"server"` // FQDN or IP address of server
	Path   string   `yaml:"path" json:"path"`     // path to store directory
	Days   int      `yaml:"days" json:"days"`
	Jobs   []string `yaml:"jobs" json:"jobs"`
	Step   int      `yaml:"step" json:"step"`
}

var (
	showConfig    = kingpin.Flag("config.show", "Show actual configuration and ends").Default("false").Bool()
	configFile    = kingpin.Flag("config.file", "Configuration file default is \"cfg.yml\".").PlaceHolder("cfg.yml").Default("cfg.yml").String()
	directoryData = kingpin.Flag("path", "Path where store export json data").PlaceHolder("path").Default("./dump").String()
	server        = kingpin.Flag("server", "Prometheus server FQDN or IP address").PlaceHolder("server").Default("").String()
	config        = &configuration{
		Server: "",
		Path:   "./dump",
		Days:   1,
		Jobs:   []string{},
		Step:   10,
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

func dirAccessible(directory string) bool {
	var file string
	dir, err := filepath.Abs(directory)
	if err != nil {
		dir = directory
	}
	for ok := true; ok; ok = fileExists(file) {
		file = filepath.Join(dir, randomString()+".tmp")
	}
	f, err := os.Create(file)
	if err == nil {
		_ = f.Close()
		e := os.Remove(file)
		if e != nil {
			fmt.Printf("Problem test access to store directory. " + e.Error())
			os.Exit(1)
		}
	}
	return err == nil
}

func (c *configuration) overWriteFromLine() {
	if len(*server) > 0 {
		c.Server = *server
	}
	if len(*directoryData) > 0 {
		c.Path = *directoryData
	}
}

func (c *configuration) validate() error {
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
	p, err := filepath.Abs(c.Path)
	if err == nil {
		c.Path = p
	}
	if c.Days < 1 || c.Days > 60 {
		return errors.New("defined days back not valid (1 - 60)")
	}
	if c.Step < 5 || c.Step > 3600 {
		c.Step = 10
	}
	if !dirExists(c.Path) {
		return errors.New("path not exists")
	}
	if !dirAccessible(c.Path) {
		return errors.New("path not accessible for write")
	}
	return nil
}

func (c *configuration) loadFile(filename string) error {
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
	c.overWriteFromLine()
	return c.validate()
}

func (c *configuration) print() string {
	a := fmt.Sprintf("\r\n%s\r\nActual configuration:\r\n", applicationName)
	a = fmt.Sprintf("%sServer:       [%s]\r\n", a, c.Server)
	a = fmt.Sprintf("%sData path:    [%s]\r\n", a, c.Path)
	a = fmt.Sprintf("%sDays back:    [%d]\r\n", a, c.Days)
	if len(c.Jobs) == 0 {
		a = fmt.Sprintf("%sTargets:      [--all--]\r\n", a)
	} else {
		a = fmt.Sprintf("%sJobs:         [%s]\r\n", a, strings.Join(c.Jobs, ", "))
	}
	return a
}

func (c *configuration) filePath(fileName string) string {
	return filepath.Join(c.Path, fileName)
}
