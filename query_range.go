package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log/level"
	"time"
)

type Result struct {
	ResultType string `json:"resultType"`
	Result     Matrix `json:"result"`
}

func UnmarshalResult(data []byte) (*Result, error) {
	var r Result
	err := json.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func getRangeDay(job string, dayBack int) (*Matrix, error) {
	start := time.Now().UTC().Add(time.Duration(-24*dayBack) * time.Hour).Truncate(time.Hour * 24)
	end := start.Add(24 * time.Hour)
	uri := fmt.Sprintf("query_range?query=%s&start=%s&end=%s&step=%d", job, start.Format(time.RFC3339), end.Format(time.RFC3339), config.Step)

	read, err := GetApiData(uri)
	if err != nil {
		_ = level.Error(logger).Log("msg", "problem collect data for job")
		return nil, err
	}
	if !read.statusSuccess() {
		_ = level.Error(logger).Log("msg", "target meta data read return error", "error", read.Error, "errorType", read.ErrorType)
		return nil, errors.New(read.Error)
	}
	s := string(*read.Data)
	t, err := UnmarshalResult([]byte(s))
	if err != nil {
		return nil, err
	}
	return &t.Result, nil
}
