package bootstrap

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	. "deblasis.net/space-traffic-control/common/utils"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/kong/go-kong/kong"
)

const (
	ProtobufBasePath = "/proto/"
)

type apiGatewayConfigurator struct {
	k      *kong.Client
	logger log.Logger
}

func NewApiGatewayConfigurator(baseUrl string, logger log.Logger) *apiGatewayConfigurator {
	kongClient, err := kong.NewClient(&baseUrl, nil)
	if err != nil {
		panic(err)
	}
	return &apiGatewayConfigurator{
		k:      kongClient,
		logger: log.With(logger, "component", "apiGatewayConfigurator"),
	}
}

func (a *apiGatewayConfigurator) RegisterGRPCService(serviceName, host, port string, tags ...string) (*kong.Service, error) {
	level.Info(a.logger).Log(
		"operation", "RegisterGRPCService",
		"msg", fmt.Sprintf("registering %v to %v:%v", serviceName, host, port),
		"tags", tags,
	)

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	service, err := a.k.Services.Create(context.TODO(), &kong.Service{
		Name:     &serviceName,
		Protocol: String("grpc"),
		Host:     &host,
		Port:     Int(portInt),
		Tags:     StringSlice(tags...),
	})
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (a *apiGatewayConfigurator) GetGRPCService(serviceName string) *kong.Service {
	level.Info(a.logger).Log(
		"operation", "GetGRPCService",
		"msg", fmt.Sprintf("checking if %v already exists and returning it", serviceName),
	)
	service, err := a.k.Services.Get(context.TODO(), String(serviceName))
	if err != nil {
		return nil
	}
	return service
}
func (a *apiGatewayConfigurator) GetHTTPRoute(serviceID string, paths []string) *kong.Route {
	level.Info(a.logger).Log(
		"operation", "GetHTTPRoute",
		"msg", fmt.Sprintf("checking if %v already exists and returning it", GetRouteName(serviceID, paths)),
	)
	route, err := a.k.Routes.Get(context.TODO(), String(GetRouteName(serviceID, paths)))
	if err != nil {
		return nil
	}
	return route
}

func (a *apiGatewayConfigurator) GetGRPCGatewayPlugin(routeNameOrID string) *kong.Plugin {
	level.Info(a.logger).Log(
		"operation", "GetGRPCGatewayPlugin",
		"msg", fmt.Sprintf("checking if %v already exists and returning it", routeNameOrID),
	)
	route, err := a.k.Routes.Get(context.TODO(), String(routeNameOrID))
	if err != nil {
		return nil
	}

	plugins, err := a.k.Plugins.ListAllForRoute(context.TODO(), route.ID)
	if err != nil {
		return nil
	}

	for _, p := range plugins {
		if p.Name == String("grpc-gateway") {
			return p
		}
	}

	return nil
}

func (a *apiGatewayConfigurator) AddHTTPRoute(serviceID *string, paths []string, tags ...string) (*kong.Route, error) {
	level.Info(a.logger).Log(
		"operation", "AddHTTPRoute",
		"msg", fmt.Sprintf("registering http route for service %v and naming it %v", *serviceID, GetRouteName(*serviceID, paths)),
		"tags", tags,
	)
	route, err := a.k.Routes.CreateInService(context.TODO(), serviceID, &kong.Route{
		Name:      String(GetRouteName(*serviceID, paths)),
		Protocols: StringSlice("http"),
		Paths:     StringSlice(paths...),
		Tags:      StringSlice(tags...),
	})
	if err != nil {
		return nil, err
	}
	return route, nil
}

func (a *apiGatewayConfigurator) AddRouteGRPCGatewayPlugin(route *kong.Route, protoPath string, tags ...string) (*kong.Plugin, error) {
	level.Info(a.logger).Log(
		"operation", "AddRouteGRPCGatewayPlugin",
		"msg", fmt.Sprintf("registering grpc-gateway plugin to route %v with protopath %v", route.ID, protoPath),
		"tags", tags,
	)
	plugin, err := a.k.Plugins.Create(context.TODO(), &kong.Plugin{
		Name:   String("grpc-gateway"),
		Route:  route,
		Config: map[string]interface{}{"proto": ProtobufBasePath + protoPath},
		Tags:   StringSlice(tags...),
	})
	if err != nil {
		return nil, err
	}
	return plugin, nil
}

func (a *apiGatewayConfigurator) DeregisterGRPCService(serviceName string) error {
	level.Info(a.logger).Log(
		"operation", "DeregisterGRPCService",
		"msg", fmt.Sprintf("deregistering service %v", serviceName),
	)

	svc := a.GetGRPCService(serviceName)
	if svc == nil {
		return nil
	}

	routes, _, _ := a.k.Routes.ListForService(context.TODO(), svc.ID, nil)
	for _, route := range routes {
		plugins, _ := a.k.Plugins.ListAllForRoute(context.TODO(), route.ID)
		for _, p := range plugins {
			a.k.Plugins.Delete(context.TODO(), p.ID)
		}
		a.k.Routes.Delete(context.TODO(), route.ID)
	}

	return a.k.Services.Delete(context.TODO(), String(serviceName))
}

func GetRouteName(serviceName string, paths []string) string {

	ret := serviceName + "~"

	for _, s := range paths {
		ret += strings.Replace(s, "/", "_", -1) + "."
	}

	return ret
}
