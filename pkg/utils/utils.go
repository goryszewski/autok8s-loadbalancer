package utils

import (
	"os/exec"

	"github.com/goryszewski/libvirtApi-client/libvirtApiClient"
)

func Add(ip string) error {

	cmd := exec.Command("/usr/sbin/ip", "addr", "add", ip, "dev", "lo")
	err := cmd.Run()

	if err != nil {
		return err
	}
	return err

}

func Del(ip string) error {

	cmd := exec.Command("/usr/sbin/ip", "addr", "del", ip+"/32", "dev", "lo")
	err := cmd.Run()

	if err != nil {
		return err
	}
	return err
}

func Compare(arr1 []libvirtApiClient.LoadBalancer, arr2 []libvirtApiClient.LoadBalancer) []libvirtApiClient.LoadBalancer {
	var diff []libvirtApiClient.LoadBalancer
	for _, lb1 := range arr1 {
		test := false
		for _, lb2 := range arr2 {
			if lb1.Namespace == lb2.Namespace && lb1.Name == lb2.Name {
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
