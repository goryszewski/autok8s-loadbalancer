package main

import (
	"fmt"
	"os"
	"text/template"
	"time"
)

type operator struct {
}
type Docker struct {
}
type LB struct {
	Name      string
	Namespace string
	Address   string
	AddressS  string
	Port      int
	PortS     int
}

func (c *operator) Add(arr2 []LB) {

	var tmplFile = "templatehaporxy.tpl"
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}
	var f *os.File
	f, err = os.Create("haproxy.cnf")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, arr2)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
}

func (c *operator) Delete(arr2 []LB) {

}

func (c *operator) GetAllDockers() []Docker {
	return nil
}

func compare(arr1 []LB, arr2 []LB) []LB {
	var diff []LB
	for _, lb1 := range arr1 {
		test := false
		for _, lb2 := range arr2 {
			if lb1.Namespace == lb2.Namespace {
				test = true
				break
			}
		}
		if !test {
			diff = append(diff, lb1)
		}
	}
	return diff
}

func main() {
	// User := "r"
	// Pass := "r"
	// URL := "http://127.0.0.1:8050"
	// conf := libvirtApiClient.Config{Username: &User, Password: &Pass, Url: &URL}
	// cc, err := libvirtApiClient.NewClient(conf, &http.Client{Timeout: 10 * time.Second})

	operator := operator{}

	// if err != nil {
	// 	fmt.Printf("ERROR:[%+v]", err)
	// }

	for true {
		fmt.Println("Start loop")
		// get all loadbalancer
		// lbs, err := cc.GetLoadBalancerAll()
		// if err != nil {
		// 	fmt.Printf("[Error] [%v]", err)
		// }

		// dockerlbs, err := operator.GetAllDockers()

		// lb_to_delete := compare(lbs, dockerlbs)

		// lb_to_add := compare(dockerlbs, lbs)

		// operator.Delete(lb_to_delete)
		lb_to_add := []LB{
			LB{Name: "test1", Namespace: "test1", Address: "10.10.10.1", Port: 80, AddressS: "10.17.3.1", PortS: 22},
			LB{Name: "test2", Namespace: "test2", Address: "10.10.10.2", Port: 80, AddressS: "10.17.3.1", PortS: 22},
		}
		operator.Add(lb_to_add)

		fmt.Println("End loop")
		time.Sleep(time.Second * 10)
	}

}
