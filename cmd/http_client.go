package cmd

import (
	"net/http"
	"time"
)

func HTTPClient() *http.Client {
	var client = &http.Client{
		Timeout: 60 * time.Second,
	}
	return client
}
