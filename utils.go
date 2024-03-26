package main

import "github.com/goryszewski/libvirtApi-client/libvirtApiClient"

func compare(arr1 []libvirtApiClient.ServiceLoadBalancerResponse, arr2 []libvirtApiClient.ServiceLoadBalancerResponse) []libvirtApiClient.ServiceLoadBalancerResponse {
	var diff []libvirtApiClient.ServiceLoadBalancerResponse
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
