package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/goryszewski/libvirtApi-client/libvirtApiClient"
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
}

type Container struct {
	name string
	id   string
}

func (d *Docker) CreateAndStart(loadbalancer libvirtApiClient.ServiceLoadBalancerResponse) error {
	err := d.Create(loadbalancer)
	if err != nil {
		return err
	}

	err = d.Start(loadbalancer)
	if err != nil {
		return err
	}

	return nil
}

func (d *Docker) Create(loadbalancer libvirtApiClient.ServiceLoadBalancerResponse) error {
	url := "http://127.0.0.1:5555" // TODO load from config file
	file_name := "/tmp/" + loadbalancer.Name + "_" + loadbalancer.Namespace + "_" + "haproxy.cnf"
	name := loadbalancer.Namespace + "_" + loadbalancer.Name

	prep_ExposedPorts := make(map[string]struct{})
	prep_PortBindings := make(map[string][]host_bind)
	for _, port := range loadbalancer.Ports {
		fmt.Printf("%v/%v \n", port.Port, port.Protocol)
		prep_ExposedPorts[fmt.Sprintf("%v/%v", port.Port, port.Protocol)] = struct{}{}
		prep_PortBindings[fmt.Sprintf("%v/%v", port.Port, port.Protocol)] = append([]host_bind{}, host_bind{HostIp: loadbalancer.Ip, HostPort: fmt.Sprintf("%v", port.Port)})

	}
	HostConfig := HostConfig{Binds: []string{file_name + ":/usr/local/etc/haproxy/haproxy.cfg"}, PortBindings: prep_PortBindings}

	payload := DockerRequest{
		Image: "haproxy",
		Labels: map[string]string{
			"haproxy":   "lb",
			"name":      loadbalancer.Name,
			"namespace": loadbalancer.Namespace,
			"ip":        loadbalancer.Ip},
		ExposedPorts: prep_ExposedPorts,
		HostConfig:   HostConfig,
	}
	fmt.Printf("%+v \n", payload)
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
func (d *Docker) Start(loadbalancer libvirtApiClient.ServiceLoadBalancerResponse) error {
	name := loadbalancer.Namespace + "_" + loadbalancer.Name
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
func (d *Docker) Delete(loadbalancer libvirtApiClient.ServiceLoadBalancerResponse) error {
	url := "http://127.0.0.1:5555"
	name := loadbalancer.Namespace + "_" + loadbalancer.Name
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%v/containers/%v?force=true", url, name), nil)
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

func (d *Docker) GetContainersByLabels(label string) ([]libvirtApiClient.ServiceLoadBalancerResponse, error) {

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

	docker_response := []DockerRequest{}

	err = json.Unmarshal(body, &docker_response)
	if err != nil {
		fmt.Println(err)
	}
	var output []libvirtApiClient.ServiceLoadBalancerResponse // TODO refactor type struct
	for _, item := range docker_response {
		ServiceLoadBalancer1 := libvirtApiClient.ServiceLoadBalancer{Name: item.Labels["name"], Namespace: item.Labels["namespace"]}
		tmp := libvirtApiClient.ServiceLoadBalancerResponse{ID: "2", Ip: item.Labels["Ip"], ServiceLoadBalancer: &ServiceLoadBalancer1}
		tmp.ServiceLoadBalancer.Name = item.Labels["name"]
		tmp.ServiceLoadBalancer.Namespace = item.Labels["namespace"]
		// tmp.Namespace = item.Labels["namespace"]
		output = append(output, tmp)
	}
	fmt.Println(output)
	return output, nil
}
