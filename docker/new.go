package docker

import (
	"net/http"
	"time"
)

func NewDockerOperator(URL string) Docker {
	return Docker{
		client: &http.Client{Timeout: 10 * time.Second},
		URL:    URL,
	}
}
