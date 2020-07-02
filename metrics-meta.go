package main

import (
	"encoding/json"
	"github.com/go-kit/kit/log/level"
	"os"
)

const metricsMetaFileName = "metrics-meta.json"

type metricMeta struct {
	Job    string `json:"job,omitempty"`
	Metric string `json:"metric,omitempty"`
	Type   string `json:"type,omitempty"`
	Help   string `json:"help,omitempty"`
	Unit   string `json:"unit,omitempty"`
}

type metricsMetaList []metricMeta

func newMetrics(target targetData) *metricMeta {
	return &metricMeta{
		Job:    target.Target["job"],
		Metric: target.Metric,
		Type:   target.Type,
		Help:   target.Help,
		Unit:   target.Unit,
	}
}

//func (m *metricMeta) isSame(meta metricMeta) bool {
//	return m.Job == meta.Job && m.Metric == meta.Metric
//}

func (m *metricMeta) isSameTarget(target targetData) bool {
	return m.Job == target.Target["job"] && m.Metric == target.Metric
}

func newMetricsMetaList(t targetList) *metricsMetaList {
	m := metricsMetaList{}
	for _, targetData := range t {
		if !m.exitInList(targetData) {
			m = append(m, *newMetrics(targetData))
		}
	}
	return &m
}

func (m *metricsMetaList) exitInList(data targetData) bool {
	for i := 0; i < len(*m); i++ {
		if (*m)[i].isSameTarget(data) {
			return true
		}
	}
	return false
}

func (m *metricsMetaList) onlyForJobs(jobNames []string) {
	if len(jobNames) == 0 {
		return
	}
	for i := len(*m) - 1; i >= 0; i-- {
		if !containsString(jobNames, (*m)[i].Job) {
			*m = append((*m)[:i], (*m)[i+1:]...)
		}
	}
}

func (m *metricsMetaList) saveList() {
	name := config.filePath(metricsMetaFileName)
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
