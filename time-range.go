package main

import (
	"fmt"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"time"
)

// initRange generate Range for data based on config value Days back
func initRange() v1.Range {
	return initRangeFromTo(time.Now().AddDate(0, 0, -config.Days).Truncate(time.Hour*24), time.Now())
}

// initRangeFromTo generate range for limit data not for day back
func initRangeFromTo(from time.Time, to time.Time) v1.Range {
	return v1.Range{
		Start: from.UTC().Add(time.Minute * -timeRangeOverSize).Truncate(time.Minute),
		End:   to.UTC().Add(time.Minute * timeRangeOverSize).Truncate(time.Minute),
		Step:  time.Duration(config.Step) * time.Second,
	}
}

// generateRangeTable Split required range to range array
func generateRangeTable(r v1.Range) []v1.Range {
	var ret []v1.Range
	startTime := r.Start
	endTime := startTime
	for r.End.After(endTime.Add(2 * time.Second)) {
		endTime = startTime.Add(6 * time.Hour)
		if endTime.After(r.End) {
			endTime = r.End
		}
		r := v1.Range{
			Start: startTime,
			End:   endTime,
			Step:  time.Duration(config.Step) * time.Second,
		}
		ret = append(ret, r)
		startTime = endTime
	}
	return ret
}

func printTimeRanges(r []v1.Range) string {

	a := fmt.Sprintln("Range data list:")
	for _, v := range r {
		a = fmt.Sprintf("%s              [%s - %s / %d second]\r\n", a, v.Start.Format(time.RFC3339), v.End.Format(time.RFC3339), v.Step/1000000000)
	}
	return a
}
