package main

import "github.com/prometheus/common/model"

const (
	LabelName     = "__name__"
	LabelInstance = "instance"
)

type OrganizedServices []OrganizedSeries

type OrganizedSeries struct {
	Name     string           `yaml:"name" json:"name"`
	Instance string           `yaml:"instance" json:"instance"`
	Set      []model.LabelSet `yaml:"set" json:"set"`
}

func splitToSeriesNameAndInstance(set []model.LabelSet) OrganizedServices {
	series := OrganizedServices{}

	for _, labelSet := range set {
		i := series.index(labelSet)
		if i > -1 {
			series[i].Set = append(series[i].Set, labelSet)
		} else {
			series = append(series, OrganizedSeries{
				Name:     string(labelSet[LabelName]),
				Instance: string(labelSet[LabelInstance]),
				Set:      []model.LabelSet{labelSet},
			})
		}
	}
	return series
}

func (s OrganizedServices) index(set model.LabelSet) int {
	for i, series := range s {
		if series.Name == string(set[LabelName]) && series.Instance == string(set[LabelInstance]) {
			return i
		}
	}
	return -1
}
