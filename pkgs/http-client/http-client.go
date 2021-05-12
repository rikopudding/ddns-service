package httpclient

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

var client = http.DefaultClient

func SetClient(c *http.Client) {
	client = c
}

func Get(url string, headers map[string]string) ([]byte, error) {
	return Request("GET", url, headers, nil)
}

func Post(url string, headers map[string]string, reqBody []byte) ([]byte, error) {
	return Request("POST", url, headers, reqBody)
}

func PostJson(url string, headers map[string]string, reqBody []byte) ([]byte, error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	if ct, ok := headers["Content-Type"]; !ok || ct != "application/json" {
		headers["Content-Type"] = "application/json"
	}
	return Request("POST", url, headers, reqBody)
}

func PutJson(url string, headers map[string]string, reqBody []byte) ([]byte, error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	if ct, ok := headers["Content-Type"]; !ok || ct != "application/json" {
		headers["Content-Type"] = "application/json"
	}
	return Request("PUT", url, headers, reqBody)
}

func Request(method, url string, headers map[string]string, reqBody []byte) (respBody []byte, err error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(reqBody))
	if err != nil {
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return respBody, nil
}
