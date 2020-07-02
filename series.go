package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

type series struct {
	Metric prometheus.Labels `json:"metric"`
	Points []point           `json:"values"`
}

func (s series) String() string {
	vals := make([]string, len(s.Points))
	for i, v := range s.Points {
		vals[i] = v.String()
	}
	return fmt.Sprintf("%s =>\n%s", s.Metric, strings.Join(vals, "\n"))
}

func (s series) sameMetrics(series series) bool {
	if len(s.Metric) != len(series.Metric) {
		return false
	}
	for key, metric := range series.Metric {
		if s.Metric[key] != metric {
			return false
		}
	}
	return true
}

func (s series) forJob(jobName string) bool {
	if len(jobName) == 0 {
		return false
	}
	return s.Metric["job"] == jobName
}

func (s series) forJobs(jobs []string) bool {
	if len(jobs) == 0 {
		return true
	}
	existName := false
	for _, job := range jobs {
		existName = existName || s.forJob(job)
	}
	return existName
}
