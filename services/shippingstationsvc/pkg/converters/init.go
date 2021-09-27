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
package converters

import (
	"reflect"

	pbdb "deblasis.net/space-traffic-control/gen/proto/go/centralcommand_dbsvc/v1"
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"

	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	m "gopkg.in/jeevatkm/go-model.v1"
)

func ConfigMappers() {
	// structToStructMapping := func(inV reflect.Value) (outV reflect.Value, err error) {
	// 	outType := reflect.TypeOf(outV)
	// 	ret := reflect.New(outType).Elem().Interface()
	// 	v := inV.Interface()
	// 	errs := m.Copy(ret, v)
	// 	if len(errs) > 0 {
	// 		return reflect.Zero(inV.Type()), errs[0]
	// 	}

	// 	return reflect.ValueOf(ret), nil
	// }

	// structToStruct := []struct {
	// 	in     interface{}
	// 	out    interface{}
	// 	oneWay bool
	// }{
	// 	{
	// 		in:  (*dtos.Ship)(nil),
	// 		out: (*dtos.Ship)(nil),
	// 	},
	// 	{
	// 		in:  (*dtos.Station)(nil),
	// 		out: (*dtos.Station)(nil),
	// 	},
	// }

	// for _, mapp := range structToStruct {
	// 	m.AddConversion(mapp.in, mapp.out, structToStructMapping)
	// 	if !mapp.oneWay {
	// 		m.AddConversion(mapp.out, mapp.in, structToStructMapping)
	// 	}
	// }

	m.AddConversion((**dtos.Dock)(nil), (**pb.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &pb.Dock{}
		v := in.Interface().(*dtos.Dock)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}
		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**pb.Station)(nil), (**dtos.Station)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &dtos.Station{}
		v := in.Interface().(*pb.Station)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}
		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**dtos.Station)(nil), (**pb.Station)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &pb.Station{}
		v := in.Interface().(*dtos.Station)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}
		return reflect.ValueOf(ret), nil
	})
	m.AddConversion((**pbdb.Station)(nil), (**dtos.Station)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &dtos.Station{}
		v := in.Interface().(*pbdb.Station)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}
		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**dtos.Station)(nil), (**pbdb.Station)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &pbdb.Station{}
		v := in.Interface().(*dtos.Station)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}
		return reflect.ValueOf(ret), nil
	})

	////
	m.AddConversion((**pb.Ship)(nil), (**dtos.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &dtos.Ship{}
		v := in.Interface().(*pb.Ship)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}
		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**dtos.Ship)(nil), (**pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &pb.Ship{}
		v := in.Interface().(*dtos.Ship)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}
		return reflect.ValueOf(ret), nil
	})
	m.AddConversion((**pbdb.Ship)(nil), (**dtos.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &dtos.Ship{}
		v := in.Interface().(*pbdb.Ship)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}
		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**dtos.Ship)(nil), (**pbdb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &pbdb.Ship{}
		v := in.Interface().(*dtos.Ship)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}
		return reflect.ValueOf(ret), nil
	})

	/////

	m.AddConversion((*dtos.Ship)(nil), (*pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &pb.Ship{}
		v := in.Interface().(dtos.Ship)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}
		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*[]dtos.Ship)(nil), (*[]*pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := make([]*pb.Ship, 0)
		v := in.Interface().([]dtos.Ship)
		for _, x := range v {
			r := &pb.Ship{}
			errs := m.Copy(r, x)
			if len(errs) > 0 {
				return reflect.Zero(in.Type()), errs[0]
			}
			ret = append(ret, r)
		}

		return reflect.ValueOf(ret), nil
	})
	m.AddConversion((*[]dtos.Station)(nil), (*[]*pb.Station)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := make([]*pb.Station, 0)
		v := in.Interface().([]dtos.Station)
		for _, x := range v {
			r := &pb.Station{}
			errs := m.Copy(r, x)
			if len(errs) > 0 {
				return reflect.Zero(in.Type()), errs[0]
			}
			ret = append(ret, r)
		}

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*[]*pb.Station)(nil), (*[]dtos.Station)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := make([]*dtos.Station, 0)
		v := in.Interface().([]pb.Station)
		for _, x := range v {
			r := &dtos.Station{}
			errs := m.Copy(r, x)
			if len(errs) > 0 {
				return reflect.Zero(in.Type()), errs[0]
			}
			ret = append(ret, r)
		}

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*[]*pb.Dock)(nil), (*[]*dtos.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := make([]*dtos.Dock, 0)
		v := in.Interface().([]*pb.Dock)
		for _, x := range v {
			r := &dtos.Dock{}
			errs := m.Copy(r, x)
			if len(errs) > 0 {
				return reflect.Zero(in.Type()), errs[0]
			}
			ret = append(ret, r)
		}

		return reflect.ValueOf(ret), nil
	})

}

func init() {

	//struct -> stuct
	ConfigMappers()

	// m.AddConversion((*dtos.Ship)(nil), (*dtos.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

	// 	ret := &dtos.Ship{}
	// 	v := in.Interface().(dtos.Ship)
	// 	errs := m.Copy(ret, v)
	// 	if len(errs) > 0 {
	// 		return reflect.Zero(in.Type()), errs[0]
	// 	}

	// 	return reflect.ValueOf(*ret), nil
	// })

	// m.AddConversion((*dtos.Station)(nil), (*dtos.Station)(nil), func(in reflect.Value) (reflect.Value, error) {

	// 	ret := &dtos.Station{}
	// 	v := in.Interface().(dtos.Station)
	// 	errs := m.Copy(ret, v)
	// 	if len(errs) > 0 {
	// 		return reflect.Zero(in.Type()), errs[0]
	// 	}

	// 	return reflect.ValueOf(*ret), nil
	// })
}

// 	m.AddConversion((*dtos.Ship)(nil), (**pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := &pb.Ship{}
// 		v := in.Interface().(dtos.Ship)
// 		switch v.Status {
// 		case "docked":
// 			ret.Status = pb.Ship_STATUS_DOCKED
// 		case "in-flight":
// 			ret.Status = pb.Ship_STATUS_INFLIGHT
// 		}

// 		errs := m.Copy(ret, v)
// 		if len(errs) > 0 {
// 			return reflect.Zero(in.Type()), errs[0]
// 		}

// 		return reflect.ValueOf(ret), nil
// 	})

// 	m.AddConversion((**pb.Ship)(nil), (*dtos.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := &dtos.Ship{}
// 		v := in.Interface().(*pb.Ship)
// 		errs := m.Copy(ret, v)

// 		switch v.Status {
// 		case pb.Ship_STATUS_DOCKED:
// 			ret.Status = "docked"
// 		case pb.Ship_STATUS_INFLIGHT:
// 			ret.Status = "in-flight"
// 		}

// 		if len(errs) > 0 {
// 			return reflect.Zero(in.Type()), errs[0]
// 		}

// 		return reflect.ValueOf(*ret), nil
// 	})

// 	m.AddConversion((*dtos.Station)(nil), (**pb.Station)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := &pb.Station{}
// 		v := in.Interface().(dtos.Station)
// 		errs := m.Copy(ret, v)
// 		if len(errs) > 0 {
// 			return reflect.Zero(in.Type()), errs[0]
// 		}

// 		return reflect.ValueOf(ret), nil
// 	})

// 	m.AddConversion((**pb.Station)(nil), (*dtos.Station)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := &dtos.Station{}
// 		v := in.Interface().(*pb.Station)
// 		errs := m.Copy(ret, v)
// 		if len(errs) > 0 {
// 			return reflect.Zero(in.Type()), errs[0]
// 		}

// 		return reflect.ValueOf(*ret), nil
// 	})

// 	m.AddConversion((*[]*dtos.Dock)(nil), (*[]*dtos.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := make([]*dtos.Dock, 0)
// 		v := in.Interface().([]*dtos.Dock)
// 		for _, x := range v {
// 			r := &dtos.Dock{}
// 			errs := m.Copy(r, x)
// 			if len(errs) > 0 {
// 				return reflect.Zero(in.Type()), errs[0]
// 			}
// 			ret = append(ret, r)
// 		}

// 		return reflect.ValueOf(ret), nil
// 	})

// 	m.AddConversion((*[]*dtos.Dock)(nil), (*[]*pb.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := make([]*pb.Dock, 0)
// 		v := in.Interface().([]*dtos.Dock)
// 		for _, x := range v {
// 			r := &pb.Dock{}
// 			errs := m.Copy(r, x)
// 			if len(errs) > 0 {
// 				return reflect.Zero(in.Type()), errs[0]
// 			}
// 			ret = append(ret, r)
// 		}

// 		return reflect.ValueOf(ret), nil
// 	})

// 	m.AddConversion((*[]*pb.Dock)(nil), (*[]*dtos.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := make([]*dtos.Dock, 0)
// 		v := in.Interface().([]*pb.Dock)
// 		for _, x := range v {
// 			r := &dtos.Dock{}
// 			errs := m.Copy(r, x)
// 			if len(errs) > 0 {
// 				return reflect.Zero(in.Type()), errs[0]
// 			}
// 			ret = append(ret, r)
// 		}

// 		return reflect.ValueOf(ret), nil
// 	})

// }
// }

/////////////////////////////////////////////////////////////

// 	m.AddConversion((*model.Dock)(nil), (*dtos.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := &dtos.Dock{}
// 		v := in.Interface().(model.Dock)
// 		errs := m.Copy(ret, v)
// 		if len(errs) > 0 {
// 			return reflect.Zero(in.Type()), errs[0]
// 		}

// 		return reflect.ValueOf(*ret), nil
// 	})

// 	m.AddConversion((*dtos.Dock)(nil), (*model.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := &model.Dock{}
// 		v := in.Interface().(dtos.Dock)
// 		errs := m.Copy(ret, v)
// 		if len(errs) > 0 {
// 			return reflect.Zero(in.Type()), errs[0]
// 		}

// 		return reflect.ValueOf(*ret), nil
// 	})

// 	m.AddConversion((**model.Dock)(nil), (**dtos.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := &dtos.Dock{}
// 		v := in.Interface().(*model.Dock)
// 		errs := m.Copy(ret, v)
// 		if len(errs) > 0 {
// 			return reflect.Zero(in.Type()), errs[0]
// 		}

// 		return reflect.ValueOf(ret), nil
// 	})

// 	m.AddConversion((**dtos.Dock)(nil), (**model.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := &model.Dock{}
// 		v := in.Interface().(*dtos.Dock)
// 		errs := m.Copy(ret, v)
// 		if len(errs) > 0 {
// 			return reflect.Zero(in.Type()), errs[0]
// 		}

// 		return reflect.ValueOf(ret), nil
// 	})

// 	m.AddConversion((*dtos.Station)(nil), (**pb.Station)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := &pb.Station{}
// 		v := in.Interface().(dtos.Station)
// 		errs := m.Copy(ret, v)
// 		if len(errs) > 0 {
// 			return reflect.Zero(in.Type()), errs[0]
// 		}

// 		return reflect.ValueOf(ret), nil
// 	})

// 	m.AddConversion((**dtos.Dock)(nil), (**pb.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := &pb.Dock{}
// 		v := in.Interface().(*dtos.Dock)
// 		errs := m.Copy(ret, v)
// 		if len(errs) > 0 {
// 			return reflect.Zero(in.Type()), errs[0]
// 		}

// 		return reflect.ValueOf(ret), nil
// 	})
// 	m.AddConversion((**pb.Dock)(nil), (**dtos.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := &dtos.Dock{}
// 		v := in.Interface().(*pb.Dock)
// 		errs := m.Copy(ret, v)
// 		if len(errs) > 0 {
// 			return reflect.Zero(in.Type()), errs[0]
// 		}

// 		return reflect.ValueOf(ret), nil
// 	})

// 	m.AddConversion((*dtos.Ship)(nil), (**pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := &pb.Ship{}
// 		v := in.Interface().(dtos.Ship)
// 		errs := m.Copy(ret, v)
// 		if len(errs) > 0 {
// 			return reflect.Zero(in.Type()), errs[0]
// 		}

// 		return reflect.ValueOf(ret), nil
// 	})

// 	m.AddConversion((*model.Dock)(nil), (*dtos.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := &dtos.Dock{}
// 		v := in.Interface().(model.Dock)
// 		errs := m.Copy(ret, v)
// 		if len(errs) > 0 {
// 			return reflect.Zero(in.Type()), errs[0]
// 		}

// 		return reflect.ValueOf(*ret), nil
// 	})

// 	m.AddConversion((*dtos.Dock)(nil), (*model.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := &model.Dock{}
// 		v := in.Interface().(dtos.Dock)
// 		errs := m.Copy(ret, v)
// 		if len(errs) > 0 {
// 			return reflect.Zero(in.Type()), errs[0]
// 		}

// 		return reflect.ValueOf(*ret), nil
// 	})

// 	m.AddConversion((**model.Dock)(nil), (**dtos.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := &dtos.Dock{}
// 		v := in.Interface().(*model.Dock)
// 		errs := m.Copy(ret, v)
// 		if len(errs) > 0 {
// 			return reflect.Zero(in.Type()), errs[0]
// 		}

// 		return reflect.ValueOf(ret), nil
// 	})

// 	m.AddConversion((**dtos.Dock)(nil), (**model.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := &model.Dock{}
// 		v := in.Interface().(*dtos.Dock)
// 		errs := m.Copy(ret, v)
// 		if len(errs) > 0 {
// 			return reflect.Zero(in.Type()), errs[0]
// 		}

// 		return reflect.ValueOf(ret), nil
// 	})

// 	m.AddConversion((*pb.Station)(nil), (**dtos.Station)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := &dtos.Station{}
// 		v := in.Interface().(pb.Station)
// 		errs := m.Copy(ret, v)
// 		if len(errs) > 0 {
// 			return reflect.Zero(in.Type()), errs[0]
// 		}

// 		return reflect.ValueOf(ret), nil
// 	})

// 	m.AddConversion((**dtos.Station)(nil), (**pb.Station)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := &pb.Station{}
// 		v := in.Interface().(*dtos.Station)
// 		errs := m.Copy(ret, v)
// 		if len(errs) > 0 {
// 			return reflect.Zero(in.Type()), errs[0]
// 		}

// 		return reflect.ValueOf(ret), nil
// 	})

// 	m.AddConversion((**dtos.Ship)(nil), (**pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

// 		ret := &pb.Ship{}
// 		v := in.Interface().(*dtos.Ship)
// 		errs := m.Copy(ret, v)
// 		if len(errs) > 0 {
// 			return reflect.Zero(in.Type()), errs[0]
// 		}

// 		return reflect.ValueOf(ret), nil
// 	})

// }
