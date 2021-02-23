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
	"time"
)

//const connectionTimeout = 10
const (
	TimeFormat      = "2006-01-02 15:04"
	DefaultDataPath = "./dump"
)

type Configuration struct {
	Server      string    `yaml:"server" json:"server"`           // FQDN or IP address of server
	Path        string    `yaml:"path" json:"path"`               // path to store directory
	Days        int       `yaml:"days" json:"days"`               // days back to dump
	From        string    `yaml:"from" json:"from"`               // date from in yyyy-mm-dd HH:MM format
	To          string    `yaml:"to" json:"to"`                   // date to in yyyy-mm-dd HH:MM format
	Sources     []Sources `yaml:"sources" json:"sources"`         // list of collected sources
	Step        int       `yaml:"step" json:"step"`               // step data dump 5 - 3600 sec
	StoreDirect bool      `yaml:"storeDirect" json:"storeDirect"` // don't create subdirectory with store time
}

type Sources struct {
	Instance  string `yaml:"instance" json:"instance"`   // instance names uses wildcards .+ mean all
	IncludeGo bool   `yaml:"includeGo" json:"includeGo"` // exclude standard go_ metrics (__name__)
}

var (
	showConfig    = kingpin.Flag("config.show", "Show actual configuration and ends").Short('v').Default("false").Bool()
	configFile    = kingpin.Flag("config.file", "Configuration file default is \"cfg.yml\".").Short('c').PlaceHolder("cfg.yml").Default("cfg.yml").String()
	directoryData = kingpin.Flag("path", "Path where store export json data").Short('p').PlaceHolder("path").Default(DefaultDataPath).String()
	server        = kingpin.Flag("server", "Prometheus server FQDN or IP address").Short('s').PlaceHolder("server").Default("").String()
	from          = kingpin.Flag("from", "Start datetime for export data").Short('f').PlaceHolder("yyyy-mm-dd HH:MM").Default("").String()
	to            = kingpin.Flag("to", "End datetime for export data").Short('t').PlaceHolder("yyyy-mm-dd HH:MM").Default("").String()
	back          = kingpin.Flag("back", "Export data back from now").Short('b').Default("0").Int()
	config        = &Configuration{
		Server:      "",
		Path:        DefaultDataPath,
		Days:        1,
		Sources:     []Sources{},
		Step:        10,
		From:        "",
		To:          "",
		StoreDirect: false,
	}
	configFrom time.Time = time.Date(1970, 01, 01, 0, 0, 0, 0, time.UTC)
	configTo   time.Time = time.Date(1970, 01, 01, 0, 0, 0, 0, time.UTC)
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
		file = filepath.Join(dir, RandomString()+".tmp")
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

func (c *Configuration) overWriteFromLine() {
	if len(*server) > 0 {
		c.Server = *server
	}
	if len(*directoryData) > 0 && *directoryData != DefaultDataPath {
		c.Path = *directoryData
	}
	if *back > 0 {
		c.Days = *back
	}
	if len(*from) > 0 {
		c.From = *from
	}
	if len(*to) > 0 {
		c.To = *to
	}
}

func (c *Configuration) validate() error {
	match, err := regexp.MatchString("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$", c.Server)
	if !match || err != nil {
		match, err = regexp.MatchString("^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$", c.Server)
		if !match || err != nil {
			return errors.New("defined Prometheus server address isn't valid FQDN or IP address")
		}
	}
	if len(c.Path) < 1 {
		c.Path = DefaultDataPath
	}
	p, err := filepath.Abs(c.Path)
	if err == nil {
		c.Path = p
	}
	if c.useRange() {
		rx := regexp.MustCompile(`.*(\d{4}-[01]\d-[0123]\d[ T]+[012]\d:[0-5]\d).*`)
		fMatch := rx.FindStringSubmatch(c.From)
		if len(fMatch) < 2 {
			return fmt.Errorf("input from is not in valid format, expect yyyy-mm-dd HH:MM, actual is %s", c.From)
		}
		c.From = fMatch[1]
		f, err := time.Parse(TimeFormat, c.From)
		if err != nil {
			return fmt.Errorf("input from is not in valid format, expect yyyy-mm-dd HH:MM, actual is %s", c.From)
		}
		configFrom = f
		fMatch = rx.FindStringSubmatch(c.To)
		if len(fMatch) < 2 {
			return fmt.Errorf("input to is not in valid format, expect yyyy-mm-dd HH:MM, actual is %s", c.From)
		}
		c.To = fMatch[1]
		t, err := time.Parse(TimeFormat, c.To)
		if err != nil {
			return fmt.Errorf("input to is not in valid format, expect yyyy-mm-dd HH:MM, actual is %s", c.To)
		}
		configTo = t
		if f.After(t) || f.Equal(t) {
			return fmt.Errorf("time from is after or same as to")
		}
		if t.Before(time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)) {
			return fmt.Errorf("input from and to is to old")
		}
		c.Days = 0
	} else {
		if c.Days < 1 || c.Days > 60 {
			return errors.New("defined days back not valid (1 - 60)")
		}
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

	if len(c.Sources) < 1 {
		return errors.New("not define any sources")
	}
	return nil
}

func (c *Configuration) LoadFile(filename string) error {
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

func (c *Configuration) print() string {
	a := fmt.Sprintf("\r\n%s\r\nActual configuration:\r\n", applicationName)
	a = fmt.Sprintf("%sServer:       [%s]\r\n", a, c.Server)
	a = fmt.Sprintf("%sData path:    [%s]\r\n", a, c.Path)
	if c.useRange() {
		a = fmt.Sprintf("%sFrom:         [%s]\r\n", a, configFrom.Format(time.RFC3339))
		a = fmt.Sprintf("%sTo  :         [%s]\r\n", a, configTo.Format(time.RFC3339))
	} else {
		a = fmt.Sprintf("%sDays back:    [%d]\r\n", a, c.Days)
	}
	if len(c.Sources) == 0 {
		a = fmt.Sprintf("%sSources:      [N/A]\r\n", a)
	} else {
		a = fmt.Sprintf("%sSources:\r\n", a)
		for _, source := range c.Sources {
			a = fmt.Sprintf("%s              [%s]\r\n", a, source.print())
		}
	}
	return a
}

func (c *Configuration) serverAddress() string {
	return fmt.Sprintf("http://%s:9090", c.Server)
}

func (c *Configuration) useRange() bool {
	return len(c.From) > 0 && len(c.To) > 0
}

func (s *Sources) print() string {
	a := fmt.Sprintf("%s (%t)", s.Instance, s.IncludeGo)
	return a
}
