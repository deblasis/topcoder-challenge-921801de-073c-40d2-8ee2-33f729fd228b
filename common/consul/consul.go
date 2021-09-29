// The MIT License (MIT)
//
// Copyright (c) 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
package consul

import (
	"fmt"
	"net"
	"time"

	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/hashicorp/consul/api"
)

type ConsulRegister struct {
	ConsulAddress                  string   // consul address
	ServiceName                    string   // service name
	Tags                           []string // consul tags
	ServicePort                    int      //service port
	MetricsPort                    int      //service port
	DeregisterCriticalServiceAfter time.Duration
	Interval                       time.Duration
	logger                         log.Logger
	forceLocalhost                 bool
}

func NewConsulRegister(consulAddress, serviceName string, servicePort, metricsPort int, tags []string, logger log.Logger, forceLocalhost bool) *ConsulRegister {
	return &ConsulRegister{
		ConsulAddress:                  consulAddress,
		ServiceName:                    serviceName,
		Tags:                           tags,
		ServicePort:                    servicePort,
		MetricsPort:                    metricsPort,
		DeregisterCriticalServiceAfter: time.Duration(1) * time.Minute,
		Interval:                       time.Duration(20) * time.Second,
		logger:                         logger,
		forceLocalhost:                 forceLocalhost,
	}
}

func (r *ConsulRegister) NewConsulGRPCRegister() (*consulsd.Registrar, error) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = r.ConsulAddress
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}
	client := consulsd.NewClient(consulClient)

	IP := r.getLocalIP()
	if r.forceLocalhost {
		IP = "localhost"
	}

	reg := &api.AgentServiceRegistration{
		ID:   fmt.Sprintf("%v-%v-%v", r.ServiceName, IP, r.ServicePort),
		Name: r.ServiceName,
		Tags: r.Tags,
		Port: r.ServicePort,
		Meta: map[string]string{
			"metricsport": fmt.Sprintf("%v", r.MetricsPort),
		},
		Address: IP,
		Check: &api.AgentServiceCheck{
			Interval:                       r.Interval.String(),
			GRPC:                           fmt.Sprintf("%v:%v/%v", IP, r.ServicePort, r.ServiceName),
			GRPCUseTLS:                     false,
			TLSSkipVerify:                  true,
			DeregisterCriticalServiceAfter: r.DeregisterCriticalServiceAfter.String(),
		},
	}
	return consulsd.NewRegistrar(client, reg, r.logger), nil
}

func (r *ConsulRegister) NewConsulHTTPRegister() (*consulsd.Registrar, error) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = r.ConsulAddress
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}
	client := consulsd.NewClient(consulClient)

	IP := r.getLocalIP()
	if r.forceLocalhost {
		IP = "localhost"
	}

	reg := &api.AgentServiceRegistration{
		ID:   fmt.Sprintf("%v-%v-%v", r.ServiceName, IP, r.ServicePort),
		Name: r.ServiceName,
		Tags: r.Tags,
		Port: r.ServicePort,
		Meta: map[string]string{
			"metricsport": fmt.Sprintf("%v", r.MetricsPort),
		},
		Address: IP,
		Check: &api.AgentServiceCheck{
			Interval:                       r.Interval.String(),
			HTTP:                           fmt.Sprintf("http://%v:%v/%v", IP, r.ServicePort, "health"),
			TLSSkipVerify:                  true,
			DeregisterCriticalServiceAfter: r.DeregisterCriticalServiceAfter.String(),
		},
	}
	return consulsd.NewRegistrar(client, reg, r.logger), nil
}

func (r *ConsulRegister) getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				level.Debug(r.logger).Log("detected_localip", ipnet.IP.To4().String())
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
