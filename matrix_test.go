package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"testing"
)

const testMatrixSuccess = "[{\"metric\":{\"__name__\":\"cucm_calls_active\",\"instance\":\"localhost:9718\",\"job\":\"cucm_monitor\",\"server\":\"perfcucma.perflab.zoomint.com\"},\"values\":[[1593475200,\"0\"],[1593475210,\"0\"],[1593561600,\"0\"]]},{\"metric\":{\"__name__\":\"cucm_calls_active\",\"instance\":\"localhost:9718\",\"job\":\"cucm_monitor\",\"server\":\"perfcucmb.perflab.zoomint.com\"},\"values\":[[1593475200,\"0\"],[1593475210,\"0\"],[1593561600,\"858\"]]}]"

func TestMatrix_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		m       Matrix
		args    args
		wantErr bool
	}{
		{"valid data", Matrix{}, args{[]byte(testMatrixSuccess)}, false},
		{"empty", Matrix{}, args{[]byte("")}, true},
		{"empty array", Matrix{}, args{[]byte("[{}]")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMatrix_appendSeries(t *testing.T) {
	type args struct {
		series Series
	}
	tests := []struct {
		name      string
		m         Matrix
		args      args
		newSeries bool
	}{
		{"contains", Matrix{
			Series{prometheus.Labels{"key": "data", "job": "test"}, nil},
			Series{prometheus.Labels{"job": "not there"}, nil}}, args{
			Series{prometheus.Labels{"key": "data", "job": "test"}, nil},
		}, false},
		{"not contains", Matrix{
			Series{prometheus.Labels{"key": "data", "job": "test"}, nil},
			Series{prometheus.Labels{"key": "data", "job": "not there"}, nil}}, args{
			Series{prometheus.Labels{"key": "data", "job": "missing"}, nil},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originLen := len(tt.m)
			tt.m.appendSeries(tt.args.series)
			if tt.newSeries {
				if originLen >= len(tt.m) {
					t.Errorf("appendSeries() = %d, want %d", len(tt.m), originLen)
				}
			} else {
				if originLen != len(tt.m) {
					t.Errorf("appendSeries() = %d, want %d", len(tt.m), originLen)
				}
			}
		})
	}
}

func TestMatrix_containsSeries(t *testing.T) {
	type args struct {
		series Series
	}
	tests := []struct {
		name string
		m    Matrix
		args args
		want bool
	}{
		{"empty array", Matrix{}, args{}, false},
		{"contains", Matrix{
			Series{prometheus.Labels{"key": "data", "job": "test"}, nil},
			Series{prometheus.Labels{"job": "not there"}, nil}}, args{
			Series{prometheus.Labels{"key": "data", "job": "test"}, nil},
		}, true},
		{"not contains", Matrix{
			Series{prometheus.Labels{"key": "data", "job": "test"}, nil},
			Series{prometheus.Labels{"key": "data", "job": "not there"}, nil}}, args{
			Series{prometheus.Labels{"key": "data", "job": "missing"}, nil},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.containsSeries(tt.args.series); got != tt.want {
				t.Errorf("containsSeries() = %v, want %v", got, tt.want)
			}
		})
	}
}