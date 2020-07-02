package main

import (
	"encoding/json"
	"github.com/go-kit/kit/log/level"
	"os"
	"strings"
)

type Matrix []Series

func (m Matrix) String() string {
	strs := make([]string, len(m))
	for i, ss := range m {
		strs[i] = ss.String()
	}
	return strings.Join(strs, "\n")
}

func (m *Matrix) UnmarshalJSON(data []byte) error {
	var s []Series
	s = []Series{}
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*m = s
	return nil
}

func (m *Matrix) containsSeries(series Series) bool {
	if len(series.Metric) == 0 {
		return false
	}
	contain := false
	for _, metric := range *m {
		contain = contain || metric.sameMetrics(series)
	}
	return contain
}

func (m *Matrix) save(metricsName string) {
	name := config.filePath(metricsName + ".json")
	f, err := os.Create(name)
	if err != nil {
		_ = level.Error(logger).Log("msg", "problem create meta data file ", "file", name, "error", err)
		return
	}
	defer func() { _ = f.Close() }()
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		_ = level.Error(logger).Log("msg", "problem create prepare target meta data", "file", name, "error", err)
		return
	}
	_, err = f.Write(data)
	if err != nil {
		_ = level.Error(logger).Log("msg", "problem write data to file ", "file", name, "error", err)
		return
	}
	_ = level.Debug(logger).Log("msg", "target meta data success write to file ", "file", name)
}

func (m *Matrix) appendSeries(series Series) {
	if series.forJobs(config.Jobs) {
		if m.containsSeries(series) {
			for i, _ := range *m {
				if (*m)[i].sameMetrics(series) {
					(*m)[i].Points = append((*m)[i].Points, series.Points...)
				}
			}
		} else {
			*m = append(*m, series)
		}
	}
}
