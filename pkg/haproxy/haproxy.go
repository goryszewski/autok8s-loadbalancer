package haproxy

import (
	"os"
	"text/template"

	"github.com/goryszewski/libvirtApi-client/libvirtApiClient"
)

func (c *Haproxy) CreateHaproxyConfig(svc_lb libvirtApiClient.LoadBalancer) (string, error) {

	file_name := c.Path + "/" + svc_lb.Name + "_" + svc_lb.Namespace + "_" + "haproxy.cnf"

	var tmplFile = "haproxy.tpl"
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		return "", err
	}

	var f *os.File
	f, err = os.Create(file_name)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(f, svc_lb)
	if err != nil {

		return "", err
	}
	err = f.Close()
	if err != nil {
		return "", err
	}
	return file_name, nil
}

func (c *Haproxy) DeleteHaproxyConfig(svc_lb libvirtApiClient.LoadBalancer) error {
	file_name := c.Path + "/" + svc_lb.Namespace + "_" + svc_lb.Name + "_" + "haproxy.cnf"

	err := os.Remove(file_name)
	if err != nil {
		return err
	}

	return nil
}
