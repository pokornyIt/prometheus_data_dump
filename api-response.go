package main

import (
	"encoding/json"
	"github.com/go-kit/kit/log/level"
)

type errorType string

type ApiResponse struct {
	Status    string           `json:"status"`
	Data      *json.RawMessage `json:"data,omitempty"`
	ErrorType errorType        `json:"errorType,omitempty"`
	Error     string           `json:"error,omitempty"`
	Warnings  []string         `json:"warnings,omitempty"`
	Response  []byte
}

func GetApiData(uri string) (*ApiResponse, error) {
	data, err := getFormUri(uri)
	if err != nil {
		_ = level.Error(logger).Log("msg", "problem collect targets details")
		return nil, err
	}
	var api ApiResponse
	if err := json.Unmarshal(data, &api); err != nil {
		return nil, err
	}
	api.Response = make([]byte, len(data))
	copy(api.Response, data)
	return &api, nil
}

func (a *ApiResponse) statusSuccess() bool {
	return a.Status == "success"
}
