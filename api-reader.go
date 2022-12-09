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

func collectSeriesList(v1api v1.API, sources Sources, dateRange v1.Range) (labels []model.LabelSet, err error) {
	_ = level.Debug(logger).Log("msg", fmt.Sprintf("entry collect series data for instance %s", sources.Instance))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	instances := fmt.Sprintf("{%s=~\"%s\"}", LabelInstance, sources.Instance)

	_ = level.Debug(logger).Log("msg", fmt.Sprintf("instance filter: %s", instances))

	dataSet, warnings, err := v1api.Series(ctx, []string{instances}, dateRange.Start, dateRange.End)
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

func collectLabelsSeriesList(v1api v1.API, lbl []Labels, dateRange v1.Range) (labels []model.LabelSet, err error) {
	_ = level.Debug(logger).Log("msg", fmt.Sprintf("entry collect series data for instance labels"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var filter []string
	var re []*regexp.Regexp
	msg := ""
	separator := ""
	for _, l := range lbl {
		instances := fmt.Sprintf("{%s=~\"%s\"}", l.Label, l.Value)
		filter = append(filter, instances)
		_ = level.Debug(logger).Log("msg", fmt.Sprintf("instance filter: %s", instances))
		msg = fmt.Sprintf("%s%s%s=~\"%s\"", msg, separator, l.Label, l.Value)
		separator = ", "
		if len(l.ExcludeMetrics) > 0 {
			r, err := regexp.Compile(l.ExcludeMetrics)
			if err == nil {
				re = append(re, r)
			}
		}
	}

	dataSet, warnings, err := v1api.Series(ctx, filter, dateRange.Start, dateRange.End)
	if err != nil {
		_ = level.Error(logger).Log("msg", "problem query Prometheus API", "error", err)
		return nil, err
	}
	if len(warnings) > 0 {
		_ = level.Warn(logger).Log("msg", "Prometheus API return warning", "warn", err)
	}
	labels = []model.LabelSet{}
	for _, set := range dataSet {
		allowed := true
		for _, r := range re {
			if !r.Match([]byte(set[LabelName])) {
				allowed = false
				break
			}
		}
		if allowed {
			labels = append(labels, set)
		}
	}
	_ = level.Debug(logger).Log("msg", fmt.Sprintf("collect %d from %d series for instance %s", len(labels), len(dataSet), msg))
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
