package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var pointRex = regexp.MustCompile(`\[(\d+\.?\d*)\,\"([^\"]+)\"\]`)

// Point represents a single data point for a given timestamp.
type Point struct {
	T model.Time
	V float64
	U string
}

type Series struct {
	Metric prometheus.Labels `json:"metric"`
	Points []Point           `json:"values"`
}

type Matrix []Series

func (p Point) String() string {
	v := strconv.FormatFloat(p.V, 'f', -1, 64)
	return fmt.Sprintf("%v @[%v]", v, p.T)
}

type Result struct {
	ResultType string `json:"resultType"`
	Result     Matrix `json:"result"`
}

func (p *Point) UnmarshalJSON(data []byte) error {
	s := string(data)
	match := pointRex.FindStringSubmatch(s)
	if len(match) < 3 {
		return errors.New("problem parse point values")
	}
	err := json.Unmarshal([]byte(match[1]), &p.T)
	if err != nil {
		return err
	}
	p.U = p.T.Time().Format(time.RFC3339)
	p.V, err = strconv.ParseFloat(match[2], 64)
	if err != nil {
		return err
	}
	return nil
}

func (s Series) String() string {
	vals := make([]string, len(s.Points))
	for i, v := range s.Points {
		vals[i] = v.String()
	}
	return fmt.Sprintf("%s =>\n%s", s.Metric, strings.Join(vals, "\n"))
}

func (m Matrix) String() string {
	// TODO(fabxc): sort, or can we rely on order from the querier?
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
