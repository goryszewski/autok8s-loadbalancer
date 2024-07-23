package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/goryszewski/libvirtApi-client/libvirtApiClient"
)

func (d *Docker) CreateAndStart(loadbalancer libvirtApiClient.LoadBalancer, bind string) error {
	err := d.Create(loadbalancer, bind)
	if err != nil {
		return err
	}

	err = d.Start(loadbalancer)
	if err != nil {
		return err
	}

	return nil
}
func (d *Docker) Create(loadbalancer libvirtApiClient.LoadBalancer, bind string) error {

	name := loadbalancer.Namespace + "_" + loadbalancer.Name

	prep_ExposedPorts := make(map[string]struct{})
	prep_PortBindings := make(map[string][]host_bind)
	for _, port := range loadbalancer.Ports {

		prep_ExposedPorts[fmt.Sprintf("%v/%v", port.Port, port.Protocol)] = struct{}{}
		prep_PortBindings[fmt.Sprintf("%v/%v", port.Port, port.Protocol)] = append([]host_bind{}, host_bind{HostIp: loadbalancer.Ip, HostPort: fmt.Sprintf("%v", port.Port)})

	}
	HostConfig := HostConfig{Binds: []string{bind + ":/usr/local/etc/haproxy/haproxy.cfg"}, PortBindings: prep_PortBindings}

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

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%v/containers/create?name=%v", d.URL, name), bytes.NewBuffer(requestBody))
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
func (d *Docker) Start(loadbalancer libvirtApiClient.LoadBalancer) error {
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
func (d *Docker) Delete(loadbalancer libvirtApiClient.LoadBalancer) error {

	name := loadbalancer.Namespace + "_" + loadbalancer.Name
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%v/containers/%v?force=true", d.URL, name), nil)
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
func (d *Docker) GetContainersByLabels(label string) ([]libvirtApiClient.LoadBalancer, error) {

	request, err := http.NewRequest("GET", fmt.Sprintf("%v/containers/json", d.URL), nil)
	if err != nil {
		fmt.Println(err)
	}

	params := url.Values{}
	params.Add("filters", fmt.Sprintf(`{"label":["%v"]}`, label))

	request.URL.RawQuery = params.Encode()

	response, err := d.client.Do(request)
	if err != nil {
		return nil, err
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
	var output []libvirtApiClient.LoadBalancer // TODO refactor type struct
	for _, item := range docker_response {
		output = append(output, libvirtApiClient.LoadBalancer{Name: item.Labels["name"], Namespace: item.Labels["namespace"], Ip: item.Labels["Ip"]})
	}
	return output, nil
}
