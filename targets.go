package main

import (
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	//"github.com/prometheus/common"
)

type TargetData struct {
	Target prometheus.Labels `json:"target,omitempty"`
	Metric string            `json:"metric,omitempty"`
	Type   string            `json:"type,omitempty"`
	Help   string            `json:"help,omitempty"`
	Unit   string            `json:"unit,omitempty"`
}

type TargetList []TargetData

func UnmarshalTargets(data []byte) (*TargetList, error) {
	var t TargetList
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func readTargetsList() (*TargetList, error) {
	read, err := GetApiData("targets/metadata")
	if err != nil {
		_ = level.Error(logger).Log("msg", "problem collect targets details")
		return nil, err
	}
	if !read.statusSuccess() {
		_ = level.Error(logger).Log("msg", "target meta data read return error", "error", read.Error, "errorType", read.ErrorType)
		return nil, errors.New(read.Error)
	}
	s := string(*read.Data)
	t, err := UnmarshalTargets([]byte(s))
	return t, err
}

func (t *TargetList) cleanAndFilterJobs(jobNames []string) *MetricsMetaList {
	m := MetricsMetaList{}
	if jobNames == nil || len(jobNames) == 0 {
		return &m
	}
	for i := len(*t) - 1; i >= 0; i-- {
		if !containsString(jobNames, (*t)[i].Target["job"]) {
			*t = append((*t)[:i], (*t)[i+1:]...)
		}
	}

	return &m
}
