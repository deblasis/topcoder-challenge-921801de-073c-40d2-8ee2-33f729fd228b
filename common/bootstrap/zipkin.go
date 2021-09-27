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
