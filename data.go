package main

import (
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/model"
	"time"
)

type SaveAllData struct {
	Collection model.LabelSet `yaml:"collection" json:"collection"`
	Data       []SaveData     `yaml:"data" json:"data"`
}

type SaveData struct {
	DateTime string  `yaml:"dateTime" json:"dateTime"`
	Value    float64 `yaml:"value" json:"value"`
}

func convertApiData(data model.Value) []SaveData {
	switch data.Type() {
	case model.ValScalar:
		_ = level.Warn(logger).Log("msg", "scalar data model not support now")
		break
	case model.ValVector:
		_ = level.Warn(logger).Log("msg", "vector data model not support now")
		break
	case model.ValMatrix:
		return dataMatrix(data)
	case model.ValString:
		_ = level.Warn(logger).Log("msg", "string data model not support now")
		break
	case model.ValNone:
		_ = level.Warn(logger).Log("msg", "collected data hasn't any data model type")
		break
	default:
		_ = level.Error(logger).Log("msg", fmt.Sprintf("unknown data model type %d", data.Type()), "error", "unknown data model type")
	}
	return []SaveData{}
}

func dataMatrix(data model.Value) []SaveData {
	matrix := data.(model.Matrix)
	var saveData []SaveData
	if matrix.Len() < 1 {
		_ = level.Error(logger).Log("msg", "in returned data is zero values")
		return nil
	}
	idx := 0
	max := len(matrix[idx].Values)
	maxLabels := len(matrix[idx].Metric)

	if matrix.Len() > 1 {
		for i, stream := range matrix {
			if len(stream.Values) > max {
				max = len(stream.Values)
				idx = i
				maxLabels = len(stream.Metric)
			} else if len(stream.Values) == max && len(stream.Metric) > maxLabels {
				maxLabels = len(stream.Metric)
				max = len(stream.Values)
				idx = i
			}
		}
		//_ = level.Debug(logger).Log("msg", "in returned data are more than one metrics, use only one")
	}

	for _, value := range matrix[idx].Values {
		saveData = append(saveData, SaveData{
			DateTime: time.Unix(0, int64(value.Timestamp)*int64(time.Millisecond)).Format(time.RFC3339),
			Value:    float64(value.Value),
		})
	}
	return saveData
}
