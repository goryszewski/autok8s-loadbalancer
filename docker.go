package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Dorequester interface {
	Do(req *http.Request) (*http.Response, error)
}

type Bind struct {
	Type        string `json"type"`        //: "bind",
	Source      string `json"source"`      //: "/home/michal/git/libvirtApi/seeder/init.sql",
	Destination string `json"destination"` //: "/docker-entrypoint-initdb.d/init.sql",
	Mode        string `json"mode"`        //: "rw",
	RW          bool   `json"rw"`
	Propagation string `json"propagation"`
}

type DockerResponse struct {
	Names []string `json:"names"`
	ID    string   `json:"id"`
	// Mount  []Bind            `json:"mounts"`
	Labels map[string]string `json"labels"`
	Image  string            `json"image"`
}

type DockerRequest struct {
	Image  string            `json"image"`
	Labels map[string]string `json"labels"`
}

type Docker struct {
	client Dorequester
}

type Container struct {
	name string
	id   string
}

func (d *Docker) Create(name string) error {
	url := "http://127.0.0.1:5555"

	payload := DockerRequest{Image: "nginx", Labels: map[string]string{"haproxy": "123"}}
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%v/containers/create?name=%v", url, name), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := d.client.Do(request)
	if err != nil {
		return err
	}

	fmt.Println(response.StatusCode)

	defer response.Body.Close()

	return nil

}
func (d *Docker) Start(name string) error {
	url := "http://127.0.0.1:5555"
	request1, err := http.NewRequest("POST", fmt.Sprintf("%v/containers/%v/start", url, name), nil)
	if err != nil {
		return err
	}
	request1.Header.Set("Content-Type", "application/json")
	response2, err := d.client.Do(request1)
	if err != nil {
		return err
	}
	fmt.Println(response2.StatusCode)

	defer response2.Body.Close()
	return nil
}
func (d *Docker) Delete(id string) error {
	url := "http://127.0.0.1:5555"
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%v/containers/%v?force=true", url, id), nil)
	if err != nil {
		return err
	}
	response, err := d.client.Do(request)
	if err != nil {
		return err
	}

	fmt.Println(response.StatusCode)
	defer response.Body.Close()

	return nil

}

func (d *Docker) GetContainersByLabels(label string) ([]DockerResponse, error) {

	// build request
	url := "http://127.0.0.1:5555"
	filter := "%7B%22label%22%3A%5B%22" + label + "%22%5D%7D"
	request, err := http.NewRequest("GET", fmt.Sprintf("%v/containers/json?filters=%v", url, filter), nil)
	if err != nil {
		fmt.Println(err)
	}
	response, err := d.client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	docker_response := []DockerResponse{}

	err = json.Unmarshal(body, &docker_response)
	if err != nil {
		fmt.Println(err)
	}

	return docker_response, nil
}
