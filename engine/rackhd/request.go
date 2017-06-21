package rackhd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	ErrHTTP404 = errors.New("HTTP request 404 NOT FOUND")
)

func send(method, apiURL string, data []byte) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, apiURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	log.Printf("Request %s %s\n", method, apiURL)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func get(apiURL string, v interface{}) error {
	return request("GET", apiURL, nil, v)
}

func request(method, apiURL string, data []byte, v interface{}) error {
	resp, err := send(method, apiURL, data)
	if err != nil {
		return err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode == 404 {
		return ErrHTTP404
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bs, v)
}
