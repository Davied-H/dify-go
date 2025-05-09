package dify

import "net/http"

type ClientConfig struct {
	ApiBaseUrl string
	HttpClient *http.Client
}

type Option func(*ClientConfig)

func DefaultConfig(apiUrl string) *ClientConfig {
	return &ClientConfig{
		ApiBaseUrl: apiUrl,
		HttpClient: &http.Client{},
	}
}
