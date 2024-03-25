package main

import (
	"os"
	"text/template"

	"github.com/goryszewski/libvirtApi-client/libvirtApiClient"
)

type operator struct {
}

func (c *operator) createHaproxyConfig(svc_lb libvirtApiClient.ServiceLoadBalancerResponse) error {

	var tmplFile = "templatehaporxy.tpl"
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		return err
	}
	var f *os.File
	f, err = os.Create(svc_lb.Name + "_" + svc_lb.Namespace + "_" + "haproxy.cnf")
	if err != nil {
		return err
	}
	err = tmpl.Execute(f, svc_lb)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *operator) deleteHaproxyConfig(svc_lb libvirtApiClient.ServiceLoadBalancerResponse) error {
	file_name := svc_lb.Name + "_" + svc_lb.Namespace + "_" + "haproxy.cnf"
	err := os.Remove(file_name)
	if err != nil {
		return err
	}
	return nil
}
