package client

import "errors"

func requestError(url string) error {
	return errors.New("an unknown error occurred while calling " + url)
}

func parseError(url string, body []byte) error {
	errorMsg := "an unknown error occurred while parsing response from " + url + ". Response was " + string(body)
	return errors.New(errorMsg)
}

func unknownError() error {
	return errors.New("an unknown error occurred")
}
