package client

import (
	"bytes"
	"net/http"
	"net/url"
	"strconv"
)

type HTTPRequestFactory interface {
	PostForm(Target, string, string, *url.Values) (*http.Request, error)
}

type UnauthenticatedRequestFactory struct{}

func (urf UnauthenticatedRequestFactory) PostForm(url string, data *url.Values) (*http.Request, error) {
	bodyBytes := []byte(data.Encode())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(bodyBytes)))

	return req, nil
}
