package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Requester interface {
	Get(client *http.Client, config Config, path string, query string) ([]byte, error)
	PostForm(client *http.Client, config Config, path string, query string, body map[string]string) ([]byte, error)
}

type UnauthenticatedRequester struct{}

func is2XX(status int) bool {
	if status >= 200 && status < 300 {
		return true
	}
	return false
}

func mapToURLValues(body map[string]string) url.Values {
	data := url.Values{}
	for key, val := range body {
		data.Add(key, val)
	}
	return data
}

func doAndRead(req *http.Request, client *http.Client, config Config) ([]byte, error) {
	if config.Verbose {
		logRequest(req)
	}

	resp, err := client.Do(req)
	if err != nil {
		if config.Verbose {
			fmt.Printf("%v\n\n", err)
		}

		return []byte{}, requestError(req.URL.String())
	}

	if config.Verbose {
		logResponse(resp)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if config.Verbose {
			fmt.Printf("%v\n\n", err)
		}

		return []byte{}, unknownError()
	}

	if !is2XX(resp.StatusCode) {
		return []byte{}, requestError(req.URL.String())
	}
	return bytes, nil
}

func (ug UnauthenticatedRequester) PostForm(client *http.Client, config Config, url string, body map[string]string) ([]byte, error) {
	data := mapToURLValues(body)

	req, err := UnauthenticatedRequestFactory{}.PostForm(url, &data)
	if err != nil {
		return []byte{}, err
	}
	return doAndRead(req, client, config)
}
