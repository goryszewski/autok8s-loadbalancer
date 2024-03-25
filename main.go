package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/goryszewski/libvirtApi-client/libvirtApiClient"
)

func main() {
	User := "r"
	Pass := "r"
	URL := "http://127.0.0.1:8050"
	conf := libvirtApiClient.Config{Username: &User, Password: &Pass, Url: &URL}
	cc, err := libvirtApiClient.NewClient(conf, &http.Client{Timeout: 10 * time.Second})

	operator := operator{}
	docker := Docker{client: &http.Client{Timeout: 10 * time.Second}}

	if err != nil {
		fmt.Printf("ERROR:[NewClient][%+v]", err)
	}

	for true {
		all_loadbalancer, err := cc.GetAllLoadBalancers()
		if err != nil {
			fmt.Printf("ERROR:[GetAllLoadBalancers]:[%+v]", err)
		}

		// docker.Create("test")
		docker.Start("test")

		all_containers, err := docker.GetContainersByLabels("haproxy")
		fmt.Println(all_containers)
		for _, item := range all_containers {
			fmt.Println(item.ID)
			// docker.Delete(item.ID)

		}
		if err != nil {
			fmt.Printf("ERROR:[GetAllContainers]:[%+v]", err)
		}

		fmt.Println("Start loop")

		exist_LoadBalancer := []libvirtApiClient.ServiceLoadBalancerResponse{}

		lb_to_add := compare(all_loadbalancer, exist_LoadBalancer)

		lb_to_delete := compare(exist_LoadBalancer, all_loadbalancer)

		fmt.Printf("DO TO ADD: %v \n", len(lb_to_add))

		for _, lb := range lb_to_add {
			operator.createHaproxyConfig(lb)
		}
		fmt.Printf("DO TO Delete: %v \n", len(lb_to_delete))
		for _, lb := range lb_to_delete {
			operator.deleteHaproxyConfig(lb)
		}

		fmt.Println("End loop")
		os.Exit(0)
		time.Sleep(time.Second * 10)
	}

}
