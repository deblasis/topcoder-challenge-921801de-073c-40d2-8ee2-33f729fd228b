package api

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type Route struct {
	Name        string
	Method      string
	Pattern     runtime.Pattern
	HandlerFunc runtime.HandlerFunc
}

type Routes []Route
