package main

import (
	"encoding/json"
	"github.com/go-kit/kit/log/level"
	"os"
)

const MetricsMetaFileName = "metrics-meta.json"

type MetricMeta struct {
	Job    string `json:"job,omitempty"`
	Metric string `json:"metric,omitempty"`
	Type   string `json:"type,omitempty"`
	Help   string `json:"help,omitempty"`
	Unit   string `json:"unit,omitempty"`
}

type MetricsMetaList []MetricMeta

func NewMetrics(target TargetData) *MetricMeta {
	return &MetricMeta{
		Job:    target.Target["job"],
		Metric: target.Metric,
		Type:   target.Type,
		Help:   target.Help,
		Unit:   target.Unit,
	}
}

func (m *MetricMeta) isSame(meta MetricMeta) bool {
	return m.Job == meta.Job && m.Metric == meta.Metric
}

func (m *MetricMeta) isSameTarget(target TargetData) bool {
	return m.Job == target.Target["job"] && m.Metric == target.Metric
}

func NewMetricsMetaList(t TargetList) *MetricsMetaList {
	m := MetricsMetaList{}
	for _, targetData := range t {
		if !m.exitInList(targetData) {
			m = append(m, *NewMetrics(targetData))
		}
	}
	return &m
}

func (m *MetricsMetaList) exitInList(data TargetData) bool {
	for i := 0; i < len(*m); i++ {
		if (*m)[i].isSameTarget(data) {
			return true
		}
	}
	return false
}

func (m *MetricsMetaList) onlyForJobs(jobNames []string) {
	if jobNames == nil || len(jobNames) == 0 {
		return
	}
	for i := len(*m) - 1; i >= 0; i-- {
		if !containsString(jobNames, (*m)[i].Job) {
			*m = append((*m)[:i], (*m)[i+1:]...)
		}
	}
}

func (m *MetricsMetaList) saveList() {
	name := config.filePath(MetricsMetaFileName)
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
