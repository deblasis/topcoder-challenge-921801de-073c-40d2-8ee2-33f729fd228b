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
package routing

import (
	"net/http"
	"strings"

	_ "deblasis.net/space-traffic-control/statik"
	"github.com/etherlabsio/healthcheck/v2"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(authGw http.Handler, ccGw http.Handler, ssGw http.Handler, healthchecks []healthcheck.Option) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	var routes = Routes{
		Route{
			"Index",
			"GET",
			"/",
			Index,
		},

		Route{
			"Login",
			strings.ToUpper("Post"),
			"/auth/login",
			authGw.ServeHTTP,
		},

		Route{
			"Signup",
			strings.ToUpper("Post"),
			"/user/signup",
			authGw.ServeHTTP,
		},

		Route{
			"ShipRegister",
			strings.ToUpper("Post"),
			"/centcom/ship/register",
			ccGw.ServeHTTP,
		},

		Route{
			"ShipsList",
			strings.ToUpper("Get"),
			"/centcom/ship/all",
			ccGw.ServeHTTP,
		},

		Route{
			"StationRegister",
			strings.ToUpper("Post"),
			"/centcom/station/register",
			ccGw.ServeHTTP,
		},

		Route{
			"StationsList",
			strings.ToUpper("Get"),
			"/centcom/station/all",
			ccGw.ServeHTTP,
		},

		Route{
			"ShipLand",
			strings.ToUpper("Post"),
			"/shipping-station/land",
			ssGw.ServeHTTP,
		},

		Route{
			"ShipRequestLanding",
			strings.ToUpper("Post"),
			"/shipping-station/request-landing",
			ssGw.ServeHTTP,
		},
	}

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		//handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	if len(healthchecks) > 0 {
		router.Handle("/health", healthcheck.Handler(healthchecks...))
	}

	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}
	staticServer := http.FileServer(statikFS)

	staticHandler := http.StripPrefix("/swaggerui/", staticServer)
	router.PathPrefix("/swaggerui/").Handler(staticHandler)

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/swaggerui/index.html", 302)
}
