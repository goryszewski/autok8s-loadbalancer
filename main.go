package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/goryszewski/autok8s-loadbalancer/docker"
	"github.com/goryszewski/autok8s-loadbalancer/haproxy"
	"github.com/goryszewski/autok8s-loadbalancer/utils"
	"github.com/goryszewski/libvirtApi-client/libvirtApiClient"
)

type Config struct {
	API struct {
		User *string
		Pass *string
		URL  *string
	}
	Docker struct {
		URL string
	}
	Haproxy struct {
		Path string
	}
}

func main() {
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(fmt.Sprintf("Problem with config file: %v \n", err))
	}

	var config Config

	err = json.Unmarshal(content, &config)
	if err != nil {
		panic(fmt.Sprintf("Problem with Unmarshal config file: %v \n", err))
	}

	fmt.Printf("%+v", config)

	conf := libvirtApiClient.Config{Username: config.API.User, Password: config.API.Pass, Url: config.API.URL}
	cc, err := libvirtApiClient.NewClient(conf, &http.Client{Timeout: 10 * time.Second})
	haproxyOperator := haproxy.NewHaproxyrOperator(config.Haproxy.Path)
	docker := docker.NewDockerOperator(config.Docker.URL)

	if err != nil {
		fmt.Printf("ERROR:[NewClient][%+v]", err)
	}

	for true {
		all_loadbalancer, err := cc.GetAllLoadBalancers()
		if err != nil {
			fmt.Printf("ERROR:[GetAllLoadBalancers]:[%+v]", err)
		}

		fmt.Println("Start loop")

		exist_LoadBalancer, err := docker.GetContainersByLabels("haproxy")
		if err != nil {
			fmt.Printf("ERROR:[GetContainersByLabels]:[%+v]", err)
		}

		lb_to_add := utils.Compare(all_loadbalancer, exist_LoadBalancer)

		lb_to_delete := utils.Compare(exist_LoadBalancer, all_loadbalancer)

		fmt.Printf("DO TO ADD: %v \n", len(lb_to_add))

		for _, lb := range lb_to_add {
			err := utils.Add(lb.Ip)
			if err != nil {
				fmt.Printf("Problem with interface:%v err: %v \n", lb.Ip, err)
			}
			config, err := haproxyOperator.CreateHaproxyConfig(lb)
			if err != nil {
				fmt.Printf("Error generate haproxy : %v \n", err)
				panic("e")
			}
			docker.CreateAndStart(lb, config)
		}

		fmt.Printf("DO TO Delete: %v \n", len(lb_to_delete))
		for _, lb := range lb_to_delete {
			haproxyOperator.DeleteHaproxyConfig(lb)
			docker.Delete(lb)
			err := utils.Del(lb.Ip)
			if err != nil {
				fmt.Printf("Proble delete ip: %v err: %v \n", lb.Ip, err)
			}
		}

		fmt.Println("End loop")

		time.Sleep(time.Second * 2)
	}

}
