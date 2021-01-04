package main

import (
	"encoding/json"
	"github.com/go-kit/kit/log/level"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// One Storage instance for one Source path
type Storage struct {
	MainPath   string `yaml:"mainPath" json:"mainPath"`
	TimePath   string `yaml:"timePath" json:"timePath"`
	SourcePath string `yaml:"sourcePath" json:"sourcePath"`
	Prepared   bool   `yaml:"prepared" json:"prepared"`
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
	storage.SourcePath = filepath.Join(storage.TimePath, cleanFilePathName(sources.Instance))
	storage.Prepared = false
	storage.Accessible = false
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

func (s *Storage) saveOrganized(services OrganizedServices) {
	data, err := json.Marshal(services)
	f := "organize-load.json"
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

var badCharacters = []string{
	"../",
	"<!--",
	"-->",
	"<",
	">",
	"'",
	"\"",
	"&",
	"$",
	"#",
	"{", "}", "[", "]", "=", ".+",
	";", "?", "%20", "%22", "+",
	"%3c",   // <
	"%253c", // <
	"%3e",   // >
	"",      // > -- fill in with % 0 e - without spaces in between
	"%28",   // (
	"%29",   // )
	"%2528", // (
	"%26",   // &
	"%24",   // $
	"%3f",   // ?
	"%3b",   // ;
	"%3d",   // =
}

func cleanFilePathName(path string) string {
	trimmed := strings.TrimSpace(path)
	trimmed = strings.Replace(trimmed, " ", "", -1)

	for _, badChar := range badCharacters {
		trimmed = strings.Replace(trimmed, badChar, "", -1)
	}
	if len(trimmed) < 5 {
		trimmed = RandomString()
	}
	return trimmed
}
