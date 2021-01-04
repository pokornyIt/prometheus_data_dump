package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log/level"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"regexp"
	"time"

	"github.com/prometheus/client_golang/api"
)

func prepareApi(configuration *Configuration) (v1api v1.API, err error) {
	client, err := api.NewClient(api.Config{
		Address: configuration.serverAddress(),
	})
	if err != nil {
		_ = level.Error(logger).Log("msg", "problem create Prometheus API client", "error", err)
		return nil, err
	}
	v1api = v1.NewAPI(client)
	return v1api, nil
}

func collectSeriesList(v1api v1.API, sources Sources, dayBack int) (labels []model.LabelSet, err error) {
	_ = level.Debug(logger).Log("msg", fmt.Sprintf("entry collect series data for instance %s", sources.Instance))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	instances := fmt.Sprintf("{%s=~\"%s\"}", LabelInstance, sources.Instance)

	dataSet, warnings, err := v1api.Series(ctx, []string{instances}, time.Now().Add(-time.Hour*time.Duration(24*dayBack)), time.Now())
	if err != nil {
		_ = level.Error(logger).Log("msg", "problem query Prometheus API", "error", err)
		return nil, err
	}
	if len(warnings) > 0 {
		_ = level.Warn(logger).Log("msg", "Prometheus API return warning", "warn", err)
	}
	labels = []model.LabelSet{}
	var re = regexp.MustCompile(`^go_.*`)
	for _, set := range dataSet {
		if !sources.IncludeGo {
			if re.Match([]byte(set[LabelName])) {
				continue
			}
		}
		labels = append(labels, set)
	}
	_ = level.Debug(logger).Log("msg", fmt.Sprintf("collect %d from %d series for instance %s", len(labels), len(dataSet), sources.Instance))
	return labels, nil
}

func readQueryRange(api v1.API, labelSet model.LabelSet, timeRange v1.Range) (data model.Value, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	data, warnings, err := api.QueryRange(ctx, labelSet.String(), timeRange)
	if err != nil {
		_ = level.Error(logger).Log("msg", "problem query Prometheus API", "error", err, "queryString", labelSet.String())
		return nil, err
	}
	if warnings != nil {
		_ = level.Warn(logger).Log("msg", "warning in query Prometheus API", "warn", warnings, "queryString", labelSet.String())
	}

	return data, nil
}
