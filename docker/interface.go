package docker

import "net/http"

type Dorequester interface {
	Do(req *http.Request) (*http.Response, error)
}
