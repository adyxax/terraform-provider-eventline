package evcli

import (
	"net/http"
	"time"
)

func NewHTTPClient() *http.Client {
	c := &http.Client{
		Timeout: 30 * time.Second,
	}

	return c
}
