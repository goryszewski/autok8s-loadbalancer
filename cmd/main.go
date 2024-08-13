package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/goryszewski/autok8s-loadbalancer/pkg/docker"
	"github.com/goryszewski/autok8s-loadbalancer/pkg/haproxy"
	"github.com/goryszewski/autok8s-loadbalancer/pkg/utils"
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
	// sigs := make(chan os.Signal, 1)

	// BEGIN flag
	var verbose string
	var dev bool
	flag.StringVar(&verbose, "verbose", "info", "Level log")
	flag.BoolVar(&dev, "dev", false, "Enable local develop")
	flag.Parse()
	// END flag

	// BEGIN load config
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(fmt.Sprintf("Problem with config file: %v \n", err))
	}

	var config Config

	err = json.Unmarshal(content, &config)
	if err != nil {
		panic(fmt.Sprintf("Problem with Unmarshal config file: %v \n", err))
	}
	// END load config

	conf := libvirtApiClient.Config{Username: config.API.User, Password: config.API.Pass, Url: config.API.URL}
	libvirtApi_client, err := libvirtApiClient.NewClient(conf, &http.Client{Timeout: 10 * time.Second})
	if err != nil {
		panic(fmt.Sprintf("Problem with Client libvirtApi: %v \n", err))
	}

	haproxyOperator := haproxy.NewHaproxyrOperator(config.Haproxy.Path)

	docker := docker.NewDockerOperator(config.Docker.URL)

	if err != nil {
		panic(fmt.Sprintf("ERROR:[NewClient][%+v]", err))
	}

	for {

		all_loadbalancer, err := libvirtApi_client.GetAllLoadBalancers()
		if err != nil {
			fmt.Printf("ERROR:[GetAllLoadBalancers]:[%+v]", err)
		}

		if dev == true {
			for _, ln := range all_loadbalancer {
				fmt.Printf("[]%+v \n", ln)
				config, err := haproxyOperator.CreateHaproxyConfig(ln)
				fmt.Printf("[%+v] [%v]\n", config, err)
			}
			fmt.Printf("DEVELOP: BREAK\n")
			// time.Sleep(time.Second * 10)
			break
		}

		exist_LoadBalancer, err := docker.GetContainersByLabels("haproxy")
		if err != nil {
			fmt.Printf("ERROR:[GetContainersByLabels]:[%+v]\n", err)
			time.Sleep(time.Second * 10)
			continue
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
