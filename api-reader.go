package main

import (
	"crypto/tls"
	"fmt"
	"github.com/go-kit/kit/log/level"
	"io/ioutil"
	"net/http"
	"time"
)

const UriFormat = "http://%s:9090/api/v1/%s"

func getFormUri(uri string) (data []byte, err error) {
	uri = fmt.Sprintf(UriFormat, config.Server, uri)
	_ = level.Debug(logger).Log("msg", "read data from server ", "uri", uri)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		_ = level.Error(logger).Log("msg", "problem create request for uri "+uri, "error", err)
		return nil, err
	}
	return finishApiRequest(req)
}

//func postFromUri(uri string, body []byte) (data []byte, err error) {
//	uri = fmt.Sprintf(UriFormat, config.Server, uri)
//	_ = level.Info(logger).Log("msg", "post and read data from server ", "uri", uri)
//	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
//	if err != nil {
//		_ = level.Error(logger).Log("msg", "problem create request for uri "+uri, "error", err)
//		return nil, err
//	}
//	return finishApiRequest(req)
//}

func finishApiRequest(req *http.Request) (data []byte, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: time.Duration(connectionTimeout) * time.Second, Transport: tr}
	_ = level.Debug(logger).Log("msg", "try read data from uri")
	resp, err := client.Do(req)
	if err != nil {
		_ = level.Error(logger).Log("msg", "problem get data from server", "error", err)
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	bodies, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		_ = level.Error(logger).Log("msg", "problem read data from response", "error", err)
		return nil, err
	}
	_ = level.Debug(logger).Log("msg", fmt.Sprintf("success read %d bytes from uri", len(bodies)))
	return bodies, nil
}
