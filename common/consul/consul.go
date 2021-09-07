package consul

import (
	"fmt"
	"net"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	consulsd "github.com/go-kit/kit/sd/consul"
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
			GRPCUseTLS:                     false, //TODO config
			TLSSkipVerify:                  true,  //TODO config
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
			TLSSkipVerify:                  true, //TODO config
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
