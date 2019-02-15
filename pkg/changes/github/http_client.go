package github

import (
	"net/http"
)

var httpClient = http.DefaultClient

func HttpClient() *http.Client {
	return httpClient
}
