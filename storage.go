package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log/level"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// One Storage instance for one Source path
type Storage struct {
	MainPath   string `yaml:"-" json:"-"`
	TimePath   string `yaml:"-" json:"-"`
	SourcePath string `yaml:"sourcePath" json:"sourcePath"`
	Prepared   bool   `yaml:"-" json:"-"`
	Accessible bool   `yaml:"accessible" json:"accessible"`
}

func NewStorage(path string, sources Sources) *Storage {
	storage := Storage{}
	var err error
	storage.MainPath, err = filepath.Abs(path)
	if err != nil {
		storage.MainPath = path
	}
	storage.TimePath = filepath.Join(storage.MainPath, timeStart.Format("20060102-1504"))
	if config.StoreDirect {
		storage.TimePath = storage.MainPath
	}
	storage.SourcePath = filepath.Join(storage.TimePath, cleanFilePathName(sources.Instance))
	storage.Prepared = false
	storage.Accessible = false
	_ = level.Debug(logger).Log("msg", fmt.Sprintf("final main store path %s", storage.TimePath))
	return &storage
}

func (s *Storage) prepareDirectory() error {
	if !s.Prepared {
		if !dirExists(s.SourcePath) {
			err := os.MkdirAll(s.SourcePath, os.ModePerm)
			if err != nil {
				_ = level.Error(logger).Log("msg", "problem create directory and path for store final data", "error", err, "path", s.SourcePath)
				return err
			}
		}
	}
	s.Prepared = true
	s.Accessible = dirAccessible(s.SourcePath)
	_ = level.Debug(logger).Log("msg", "directory for store final data exist and accessible", "path", s.SourcePath)
	return nil
}

func (s *Storage) saveOrganized(services OrganizedServices, timeRange v1.Range) {
	d := OrganizedForFile{
		From: timeRange.Start.Format(time.RFC3339),
		To:   timeRange.End.Format(time.RFC3339),
		Step: config.Step,
		Data: services,
	}

	data, err := json.Marshal(d)
	f := fmt.Sprintf("organize-%s.json", filepath.Base(s.SourcePath))
	if err != nil {
		_ = level.Error(logger).Log("msg", "problem convert data to JSON file", "error", err, "file", f)
		return
	}
	s.saveJson(data, filepath.Join(s.TimePath, f))
}

func (s *Storage) saveAllData(saveAllData []SaveAllData, fileName string) {
	if s.Accessible {
		data, err := json.Marshal(saveAllData)
		if err != nil {
			_ = level.Error(logger).Log("msg", "problem convert data to JSON file", "error", err, "file", fileName)
			return
		}
		s.saveJson(data, filepath.Join(s.SourcePath, fileName))
	}
}

func (s *Storage) saveJson(data []byte, fullFileName string) {
	if s.Accessible {
		err := ioutil.WriteFile(fullFileName, data, os.ModePerm)
		f := filepath.Base(fullFileName)
		if err != nil {
			_ = level.Error(logger).Log("msg", "problem save file", "error", err, "file", f)
		} else {
			_ = level.Debug(logger).Log("msg", "data successful saved to disk", "file", f)
		}
	}
}

func cleanFilePathName(path string) string {
	trimmed := strings.TrimSpace(path)
	trimmed = strings.ToLower(trimmed)

	resp := ""
	re := regexp.MustCompile(`(?i)(?i)([a-z0-9 \-_~.])`)

	for i := 0; i < len(trimmed); i++ {
		if re.MatchString(string(trimmed[i])) {
			resp = resp + string(trimmed[i])
		}
	}

	if len(resp) < 1 {
		resp = RandomString()
	}
	return strings.TrimRight(resp, ".")
}
