package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"testing"
)

func TestSeries_forJob(t *testing.T) {
	type fields struct {
		Metric prometheus.Labels
		Points []point
	}
	type args struct {
		jobName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"for job", fields{Metric: prometheus.Labels{"any": "any", "job": "test_job"}, Points: nil}, args{jobName: "test_job"}, true},
		{"not for job", fields{Metric: prometheus.Labels{"any": "any", "job": "test_job_no"}, Points: nil}, args{jobName: "test_job"}, false},
		{"label not exist", fields{Metric: prometheus.Labels{"any": "any", "job_not": "test_job_no"}, Points: nil}, args{jobName: "test_job"}, false},
		{"label not exist", fields{Metric: prometheus.Labels{"any": "any", "job_not": "test_job_no"}, Points: nil}, args{jobName: "test_job"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := series{
				Metric: tt.fields.Metric,
				Points: tt.fields.Points,
			}
			if got := s.forJob(tt.args.jobName); got != tt.want {
				t.Errorf("forJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeries_forJobs(t *testing.T) {
	type fields struct {
		Metric prometheus.Labels
		Points []point
	}
	type args struct {
		jobs []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"for job", fields{Metric: prometheus.Labels{"any": "any", "job": "test_job"}, Points: nil}, args{jobs: []string{"test_job", "second"}}, true},
		{"not for job", fields{Metric: prometheus.Labels{"any": "any", "job": "test_job_no"}, Points: nil}, args{jobs: []string{"test_job", "second"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := series{
				Metric: tt.fields.Metric,
				Points: tt.fields.Points,
			}
			if got := s.forJobs(tt.args.jobs); got != tt.want {
				t.Errorf("forJobs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeries_sameMetrics(t *testing.T) {
	type fields struct {
		Metric prometheus.Labels
		Points []point
	}
	type args struct {
		series series
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"same metrics", fields{Metric: prometheus.Labels{"key": "any", "job": "test_job"}, Points: nil},
			args{series{Metric: prometheus.Labels{"key": "any", "job": "test_job"}}}, true},
		{"different value", fields{Metric: prometheus.Labels{"key": "any", "job": "test_job_no"}, Points: nil},
			args{series{Metric: prometheus.Labels{"key": "any", "job": "test_job"}}}, false},
		{"different key", fields{Metric: prometheus.Labels{"key_diff": "any", "job": "test_job"}, Points: nil},
			args{series{Metric: prometheus.Labels{"key": "any", "job": "test_job"}}}, false},
		{"empty source", fields{Metric: prometheus.Labels{"key": "any", "job": "test_job"}, Points: nil},
			args{series{Metric: prometheus.Labels{}}}, false},
		{"empty test", fields{Metric: prometheus.Labels{}, Points: nil},
			args{series{Metric: prometheus.Labels{"key": "any", "job": "test_job"}}}, false},
		{"empty both", fields{Metric: prometheus.Labels{}, Points: nil},
			args{series{Metric: prometheus.Labels{"key": "any", "job": "test_job"}}}, false},
		{"add test", fields{Metric: prometheus.Labels{"key": "any", "job": "test_job", "add": "add"}, Points: nil},
			args{series{Metric: prometheus.Labels{"key": "any", "job": "test_job"}}}, false},
		{"add source", fields{Metric: prometheus.Labels{"key": "any", "job": "test_job"}, Points: nil},
			args{series{Metric: prometheus.Labels{"key": "any", "job": "test_job", "add": "add"}}}, false},
		{"add source", fields{Metric: prometheus.Labels{"key": "any", "job": "test_job"}, Points: nil},
			args{series{Metric: prometheus.Labels{"key": "any", "job": "test_job", "add": "add"}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := series{
				Metric: tt.fields.Metric,
				Points: tt.fields.Points,
			}
			if got := s.sameMetrics(tt.args.series); got != tt.want {
				t.Errorf("sameMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
