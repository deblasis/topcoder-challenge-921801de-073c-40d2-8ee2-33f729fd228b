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
package converters

import (
	"reflect"

	pbdb "deblasis.net/space-traffic-control/gen/proto/go/centralcommand_dbsvc/v1"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"

	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	m "gopkg.in/jeevatkm/go-model.v1"
)

func ConfigMappers() {

	m.AddConversion((**dtos.Dock)(nil), (**pb.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &pb.Dock{}
		v := in.Interface().(*dtos.Dock)
		m.Copy(ret, v)
		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**pb.Station)(nil), (**dtos.Station)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &dtos.Station{}
		v := in.Interface().(*pb.Station)
		m.Copy(ret, v)
		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**dtos.Station)(nil), (**pb.Station)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &pb.Station{}
		v := in.Interface().(*dtos.Station)
		m.Copy(ret, v)
		return reflect.ValueOf(ret), nil
	})
	m.AddConversion((**pbdb.Station)(nil), (**dtos.Station)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &dtos.Station{}
		v := in.Interface().(*pbdb.Station)
		m.Copy(ret, v)
		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**dtos.Station)(nil), (**pbdb.Station)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &pbdb.Station{}
		v := in.Interface().(*dtos.Station)
		m.Copy(ret, v)
		return reflect.ValueOf(ret), nil
	})

	////
	m.AddConversion((**pb.Ship)(nil), (**dtos.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &dtos.Ship{}
		v := in.Interface().(*pb.Ship)
		m.Copy(ret, v)
		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**dtos.Ship)(nil), (**pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &pb.Ship{}
		v := in.Interface().(*dtos.Ship)
		m.Copy(ret, v)
		return reflect.ValueOf(ret), nil
	})
	m.AddConversion((**pbdb.Ship)(nil), (**dtos.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &dtos.Ship{}
		v := in.Interface().(*pbdb.Ship)
		m.Copy(ret, v)
		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**dtos.Ship)(nil), (**pbdb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &pbdb.Ship{}
		v := in.Interface().(*dtos.Ship)
		m.Copy(ret, v)
		return reflect.ValueOf(ret), nil
	})

	/////

	m.AddConversion((*dtos.Ship)(nil), (*pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &pb.Ship{}
		v := in.Interface().(dtos.Ship)
		m.Copy(ret, v)
		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*[]dtos.Ship)(nil), (*[]*pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := make([]*pb.Ship, 0)
		v := in.Interface().([]dtos.Ship)
		for _, x := range v {
			r := &pb.Ship{}
			m.Copy(r, x)
			ret = append(ret, r)
		}

		return reflect.ValueOf(ret), nil
	})
	m.AddConversion((*[]dtos.Station)(nil), (*[]*pb.Station)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := make([]*pb.Station, 0)
		v := in.Interface().([]dtos.Station)
		for _, x := range v {
			r := &pb.Station{}
			m.Copy(r, x)
			ret = append(ret, r)
		}

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*[]*pb.Station)(nil), (*[]dtos.Station)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := make([]*dtos.Station, 0)
		v := in.Interface().([]pb.Station)
		for _, x := range v {
			r := &dtos.Station{}
			m.Copy(r, x)
			ret = append(ret, r)
		}

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*[]*pb.Dock)(nil), (*[]*dtos.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := make([]*dtos.Dock, 0)
		v := in.Interface().([]*pb.Dock)
		for _, x := range v {
			r := &dtos.Dock{}
			m.Copy(r, x)
			ret = append(ret, r)
		}

		return reflect.ValueOf(ret), nil
	})

}

func init() {
	ConfigMappers()
}
