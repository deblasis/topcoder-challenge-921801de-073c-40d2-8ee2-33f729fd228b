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

	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommand_dbsvc/v1"
	pbsvc "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	m "gopkg.in/jeevatkm/go-model.v1"
)

//TODO DRY
func init() {

	m.AddConversion((*float32)(nil), (**float32)(nil), func(in reflect.Value) (reflect.Value, error) {
		r := float32(in.Float())
		return reflect.ValueOf(&r), nil
	})

	m.AddConversion((*int64)(nil), (**int64)(nil), func(in reflect.Value) (reflect.Value, error) {
		r := in.Int()
		return reflect.ValueOf(&r), nil
	})

	m.AddConversion((*int32)(nil), (**int32)(nil), func(in reflect.Value) (reflect.Value, error) {
		r := int32(in.Int())
		return reflect.ValueOf(&r), nil
	})

	m.AddConversion((**float32)(nil), (*float32)(nil), func(in reflect.Value) (reflect.Value, error) {
		return reflect.ValueOf(float32(in.Elem().Float())), nil
	})

	m.AddConversion((**int64)(nil), (*int64)(nil), func(in reflect.Value) (reflect.Value, error) {
		return reflect.ValueOf(in.Elem().Int()), nil
	})

	m.AddConversion((**int32)(nil), (*int32)(nil), func(in reflect.Value) (reflect.Value, error) {
		return reflect.ValueOf(int32(in.Elem().Int())), nil
	})

	m.AddConversion((*model.Dock)(nil), (*dtos.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &dtos.Dock{}
		v := in.Interface().(model.Dock)
		m.Copy(ret, v)

		return reflect.ValueOf(*ret), nil
	})

	m.AddConversion((*dtos.Dock)(nil), (*model.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &model.Dock{}
		v := in.Interface().(dtos.Dock)
		m.Copy(ret, v)

		return reflect.ValueOf(*ret), nil
	})

	m.AddConversion((**model.Dock)(nil), (**dtos.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &dtos.Dock{}
		v := in.Interface().(*model.Dock)
		m.Copy(ret, v)

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**dtos.Dock)(nil), (**model.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &model.Dock{}
		v := in.Interface().(*dtos.Dock)
		m.Copy(ret, v)

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*dtos.Station)(nil), (**pb.Station)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &pb.Station{}
		v := in.Interface().(dtos.Station)
		m.Copy(ret, v)

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**dtos.Dock)(nil), (**pb.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &pb.Dock{}
		v := in.Interface().(*dtos.Dock)
		m.Copy(ret, v)

		return reflect.ValueOf(ret), nil
	})
	m.AddConversion((**pb.Dock)(nil), (**dtos.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &dtos.Dock{}
		v := in.Interface().(*pb.Dock)
		m.Copy(ret, v)

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*dtos.Ship)(nil), (**pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &pb.Ship{}
		v := in.Interface().(dtos.Ship)
		m.Copy(ret, v)

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*model.Dock)(nil), (*dtos.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &dtos.Dock{}
		v := in.Interface().(model.Dock)
		m.Copy(ret, v)

		return reflect.ValueOf(*ret), nil
	})

	m.AddConversion((*dtos.Dock)(nil), (*model.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &model.Dock{}
		v := in.Interface().(dtos.Dock)
		m.Copy(ret, v)

		return reflect.ValueOf(*ret), nil
	})

	m.AddConversion((**model.Dock)(nil), (**dtos.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &dtos.Dock{}
		v := in.Interface().(*model.Dock)
		m.Copy(ret, v)

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**dtos.Dock)(nil), (**model.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &model.Dock{}
		v := in.Interface().(*dtos.Dock)
		m.Copy(ret, v)

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*pb.Station)(nil), (**dtos.Station)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &dtos.Station{}
		v := in.Interface().(pb.Station)
		m.Copy(ret, v)

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**dtos.Station)(nil), (**pb.Station)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &pb.Station{}
		v := in.Interface().(*dtos.Station)
		m.Copy(ret, v)

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**dtos.Ship)(nil), (**pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &pb.Ship{}
		v := in.Interface().(*dtos.Ship)
		m.Copy(ret, v)

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**pb.Ship)(nil), (*dtos.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &dtos.Ship{}
		v := in.Interface().(*pb.Ship)
		m.Copy(ret, v)

		switch v.Status {
		case pb.Ship_STATUS_DOCKED:
			ret.Status = "docked"
		case pb.Ship_STATUS_INFLIGHT:
			ret.Status = "in-flight"
		}

		return reflect.ValueOf(*ret), nil
	})

	m.AddConversion((*dtos.Ship)(nil), (**pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &pb.Ship{}
		v := in.Interface().(dtos.Ship)
		m.Copy(ret, v)

		switch v.Status {
		case "docked":
			ret.Status = pb.Ship_STATUS_DOCKED
		case "in-flight":
			ret.Status = pb.Ship_STATUS_INFLIGHT
		}

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*[]*pb.Station)(nil), (*[]dtos.Station)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := make([]dtos.Station, 0)
		v := in.Interface().([]*pb.Station)
		for _, x := range v {
			r := &dtos.Station{}
			m.Copy(r, x)
			ret = append(ret, *r)
		}

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*dtos.Station)(nil), (*pbsvc.Station)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &pbsvc.Station{}
		v := in.Interface().(dtos.Station)
		m.Copy(ret, v)
		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*pbsvc.Station)(nil), (*dtos.Station)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &dtos.Station{}
		v := in.Interface().(pbsvc.Station)
		m.Copy(ret, v)
		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*dtos.Ship)(nil), (*pbsvc.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &pbsvc.Ship{}
		v := in.Interface().(dtos.Ship)
		m.Copy(ret, v)
		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*pbsvc.Ship)(nil), (*dtos.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &dtos.Ship{}
		v := in.Interface().(pbsvc.Ship)
		m.Copy(ret, v)
		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**pb.NextAvailableDockingStation)(nil), (**dtos.NextAvailableDockingStation)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &dtos.NextAvailableDockingStation{}
		v := in.Interface().(*pb.NextAvailableDockingStation)
		m.Copy(ret, v)

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**dtos.NextAvailableDockingStation)(nil), (**pb.NextAvailableDockingStation)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &pb.NextAvailableDockingStation{}
		v := in.Interface().(*dtos.NextAvailableDockingStation)
		m.Copy(ret, v)

		return reflect.ValueOf(ret), nil
	})
}
