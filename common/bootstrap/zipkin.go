//
// Copyright 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package bootstrap

import (
	"os"

	"deblasis.net/space-traffic-control/common/config"
	"github.com/go-kit/log/level"
	"github.com/lightstep/lightstep-tracer-go"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

func SetupTracers(cfg config.Config, serviceName string) (zipkinTracer *zipkin.Tracer, tracer stdopentracing.Tracer) {
	{
		if cfg.Zipkin.V2Url != "" {
			var (
				err         error
				hostPort    = "localhost:80"
				serviceName = serviceName
				reporter    = zipkinhttp.NewReporter(cfg.Zipkin.V2Url)
			)
			defer reporter.Close()
			zEP, _ := zipkin.NewEndpoint(serviceName, hostPort)
			zipkinTracer, err = zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(zEP))
			if err != nil {
				level.Error(cfg.Logger).Log("err", err)
				os.Exit(1)
			}
			if !(cfg.Zipkin.UseBridge) {
				level.Info(cfg.Logger).Log("tracer", "Zipkin", "type", "Native", "URL", cfg.Zipkin.V2Url)
			}
		}
	}

	// Determine which OpenTracing tracer to use. We'll pass the tracer to all the
	// components that use it, as a dependency.
	{
		if cfg.Zipkin.UseBridge && zipkinTracer != nil {
			level.Info(cfg.Logger).Log("tracer", "Zipkin", "type", "OpenTracing", "URL", cfg.Zipkin.V2Url)
			tracer = zipkinot.Wrap(zipkinTracer)
			zipkinTracer = nil // do not instrument with both native tracer and opentracing bridge
		} else if cfg.Zipkin.LightstepToken != "" {
			level.Info(cfg.Logger).Log("tracer", "LightStep") // probably don't want to print out the token :)
			tracer = lightstep.NewTracer(lightstep.Options{
				AccessToken: cfg.Zipkin.LightstepToken,
			})
			defer lightstep.FlushLightStepTracer(tracer)
		} else {
			tracer = stdopentracing.GlobalTracer() // no-op
		}
	}

	return
}
