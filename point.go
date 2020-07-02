package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/prometheus/common/model"
	"regexp"
	"strconv"
)

var pointRex = regexp.MustCompile(`\[\s*(\d+\.?\d*)\s*\,\s*\"([^\"]+)\"\s*\]`)

// point represents a single data point for a given timestamp.
type point struct {
	T model.Time
	V float64
	//U string
}

func (p point) String() string {
	v := strconv.FormatFloat(p.V, 'f', -1, 64)
	return fmt.Sprintf("%v @[%v]", v, p.T)
}

// Implement json Unmarshaler interface
func (p *point) UnmarshalJSON(data []byte) error {
	s := string(data)
	match := pointRex.FindStringSubmatch(s)
	if len(match) < 3 {
		return errors.New("problem parse point values")
	}
	err := json.Unmarshal([]byte(match[1]), &p.T)
	if err != nil {
		return err
	}
	//p.U = p.T.Time().Format(time.RFC3339)
	p.V, err = strconv.ParseFloat(match[2], 64)
	if err != nil {
		return err
	}
	return nil
}
