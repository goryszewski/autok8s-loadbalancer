package docker

type Bind struct {
	Type        string `json"type"`
	Source      string `json"source"`
	Destination string `json"destination"`
	Mode        string `json"mode"`
	RW          bool   `json"rw"`
	Propagation string `json"propagation"`
}

type DockerResponse struct {
	Names  []string          `json:"names"`
	ID     string            `json:"id"`
	Labels map[string]string `json"labels"`
	Image  string            `json"image"`
}
type host_bind struct {
	HostPort string
	HostIp   string
}

type HostConfig struct {
	Binds        []string `json"Binds"`
	PortBindings map[string][]host_bind
}

type DockerRequest struct {
	Image        string              `json"image"`
	Labels       map[string]string   `json"labels"`
	ExposedPorts map[string]struct{} `json"ExposedPorts"`
	HostConfig   HostConfig          `json"HostConfig"`
}

type Docker struct {
	client Dorequester
	URL    string
}

type Container struct {
	name string
	id   string
}
