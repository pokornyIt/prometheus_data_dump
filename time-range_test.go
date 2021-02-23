package main

import (
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"reflect"
	"strings"
	"testing"
	"time"
)

const layout = "2006-01-02T15:04:05.000Z"

func Test_generateRangeTable(t *testing.T) {
	type args struct {
		r    v1.Range
		step int
	}
	tests := []struct {
		name string
		args args
		want []v1.Range
	}{
		{name: "one day", args: args{
			r: v1.Range{
				Start: time.Date(2021, 1, 10, 00, 0, 0, 0, time.UTC),
				End:   time.Date(2021, 1, 11, 14, 0, 0, 0, time.UTC),
				Step:  time.Duration(15) * time.Second,
			},
			step: 15,
		}, want: []v1.Range{
			{
				Start: time.Date(2021, 1, 10, 00, 0, 0, 0, time.UTC),
				End:   time.Date(2021, 1, 10, 06, 0, 0, 0, time.UTC),
				Step:  time.Duration(15) * time.Second,
			},
			{
				Start: time.Date(2021, 1, 10, 06, 0, 0, 0, time.UTC),
				End:   time.Date(2021, 1, 10, 12, 0, 0, 0, time.UTC),
				Step:  time.Duration(15) * time.Second,
			},
			{
				Start: time.Date(2021, 1, 10, 12, 0, 0, 0, time.UTC),
				End:   time.Date(2021, 1, 10, 18, 0, 0, 0, time.UTC),
				Step:  time.Duration(15) * time.Second,
			},
			{
				Start: time.Date(2021, 1, 10, 18, 0, 0, 0, time.UTC),
				End:   time.Date(2021, 1, 11, 00, 0, 0, 0, time.UTC),
				Step:  time.Duration(15) * time.Second,
			},
			{
				Start: time.Date(2021, 1, 11, 00, 0, 0, 0, time.UTC),
				End:   time.Date(2021, 1, 11, 06, 0, 0, 0, time.UTC),
				Step:  time.Duration(15) * time.Second,
			},
			{
				Start: time.Date(2021, 1, 11, 06, 0, 0, 0, time.UTC),
				End:   time.Date(2021, 1, 11, 12, 0, 0, 0, time.UTC),
				Step:  time.Duration(15) * time.Second,
			},
			{
				Start: time.Date(2021, 1, 11, 12, 0, 0, 0, time.UTC),
				End:   time.Date(2021, 1, 11, 14, 0, 0, 0, time.UTC),
				Step:  time.Duration(15) * time.Second,
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Step = tt.args.step
			if got := generateRangeTable(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateRangeTable() = \r\n%v\r\n want\r\n%v", got, tt.want)
			}
		})
	}
}

func Test_initRange(t *testing.T) {
	tests := []struct {
		name string
		days int
		step int
		want v1.Range
	}{
		{name: "one day", days: 1, step: 10, want: v1.Range{
			Start: time.Now().UTC().AddDate(0, 0, -config.Days).Truncate(time.Hour * 24),
			End:   time.Now().UTC().Truncate(time.Minute),
			Step:  time.Duration(10) * time.Second,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Days = tt.days
			config.Step = tt.step
			if got := initRange(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initRange() = \r\n%v\r\n%v", got, tt.want)
			}
		})
	}
}

func Test_initRangeFromTo(t *testing.T) {
	type args struct {
		from time.Time
		to   time.Time
		step int
	}
	tests := []struct {
		name string
		args args
		want v1.Range
	}{
		{name: "one day", args: args{
			from: time.Date(2021, 1, 10, 10, 5, 0, 0, time.Local),
			to:   time.Date(2021, 1, 11, 15, 0, 0, 0, time.Local),
			step: 15,
		}, want: v1.Range{
			Start: time.Date(2021, 1, 10, 00, 0, 0, 0, time.UTC),
			End:   time.Date(2021, 1, 11, 14, 0, 0, 0, time.UTC),
			Step:  time.Duration(15) * time.Second,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Step = tt.args.step
			if got := initRangeFromTo(tt.args.from, tt.args.to); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initRangeFromTo() = \r\n%v\r\n%v", got, tt.want)
			}
		})
	}
}

func Test_printTimeRanges(t *testing.T) {
	type args struct {
		r []v1.Range
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "test lines", args: args{r: []v1.Range{
			{
				Start: time.Date(2021, 1, 10, 00, 0, 0, 0, time.UTC),
				End:   time.Date(2021, 1, 10, 06, 0, 0, 0, time.UTC),
				Step:  time.Duration(15) * time.Second,
			},
			{
				Start: time.Date(2021, 1, 10, 06, 0, 0, 0, time.UTC),
				End:   time.Date(2021, 1, 10, 12, 0, 0, 0, time.UTC),
				Step:  time.Duration(15) * time.Second,
			},
			{
				Start: time.Date(2021, 1, 10, 12, 0, 0, 0, time.UTC),
				End:   time.Date(2021, 1, 10, 18, 0, 0, 0, time.UTC),
				Step:  time.Duration(15) * time.Second,
			},
		},
		}, want: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := printTimeRanges(tt.args.r); strings.Count(got, "\n") != tt.want {
				t.Errorf("printTimeRanges() = %v, want %v", got, tt.want)
			}
		})
	}
}
