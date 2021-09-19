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
	m.AddConversion((*model.Dock)(nil), (*dtos.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &dtos.Dock{}
		v := in.Interface().(model.Dock)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}

		return reflect.ValueOf(*ret), nil
	})

	m.AddConversion((*dtos.Dock)(nil), (*model.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &model.Dock{}
		v := in.Interface().(dtos.Dock)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}

		return reflect.ValueOf(*ret), nil
	})

	m.AddConversion((**model.Dock)(nil), (**dtos.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &dtos.Dock{}
		v := in.Interface().(*model.Dock)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**dtos.Dock)(nil), (**model.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &model.Dock{}
		v := in.Interface().(*dtos.Dock)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*dtos.Station)(nil), (**pb.Station)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &pb.Station{}
		v := in.Interface().(dtos.Station)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**dtos.Dock)(nil), (**pb.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &pb.Dock{}
		v := in.Interface().(*dtos.Dock)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}

		return reflect.ValueOf(ret), nil
	})
	m.AddConversion((**pb.Dock)(nil), (**dtos.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &dtos.Dock{}
		v := in.Interface().(*pb.Dock)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*dtos.Ship)(nil), (**pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &pb.Ship{}
		v := in.Interface().(dtos.Ship)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*model.Dock)(nil), (*dtos.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &dtos.Dock{}
		v := in.Interface().(model.Dock)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}

		return reflect.ValueOf(*ret), nil
	})

	m.AddConversion((*dtos.Dock)(nil), (*model.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &model.Dock{}
		v := in.Interface().(dtos.Dock)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}

		return reflect.ValueOf(*ret), nil
	})

	m.AddConversion((**model.Dock)(nil), (**dtos.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &dtos.Dock{}
		v := in.Interface().(*model.Dock)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**dtos.Dock)(nil), (**model.Dock)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &model.Dock{}
		v := in.Interface().(*dtos.Dock)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*pb.Station)(nil), (**dtos.Station)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &dtos.Station{}
		v := in.Interface().(pb.Station)
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

	m.AddConversion((**dtos.Ship)(nil), (**pb.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &pb.Ship{}
		v := in.Interface().(*dtos.Ship)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**pb.Ship)(nil), (*dtos.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &dtos.Ship{}
		v := in.Interface().(*pb.Ship)
		errs := m.Copy(ret, v)

		switch v.Status {
		case pb.Ship_STATUS_DOCKED:
			ret.Status = "docked"
		case pb.Ship_STATUS_INFLIGHT:
			ret.Status = "in-flight"
		}

		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}

		return reflect.ValueOf(*ret), nil
	})

	m.AddConversion((*[]*pb.Station)(nil), (*[]dtos.Station)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := make([]dtos.Station, 0)
		v := in.Interface().([]*pb.Station)
		for _, x := range v {
			r := &dtos.Station{}
			errs := m.Copy(r, x)
			if len(errs) > 0 {
				return reflect.Zero(in.Type()), errs[0]
			}
			ret = append(ret, *r)
		}

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*dtos.Station)(nil), (*pbsvc.Station)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &pbsvc.Station{}
		v := in.Interface().(dtos.Station)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}
		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*pbsvc.Station)(nil), (*dtos.Station)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &dtos.Station{}
		v := in.Interface().(pbsvc.Station)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}
		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*dtos.Ship)(nil), (*pbsvc.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &pbsvc.Ship{}
		v := in.Interface().(dtos.Ship)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}
		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((*pbsvc.Ship)(nil), (*dtos.Ship)(nil), func(in reflect.Value) (reflect.Value, error) {
		ret := &dtos.Ship{}
		v := in.Interface().(pbsvc.Ship)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}
		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**pb.NextAvailableDockingStation)(nil), (**dtos.NextAvailableDockingStation)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &dtos.NextAvailableDockingStation{}
		v := in.Interface().(*pb.NextAvailableDockingStation)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}

		return reflect.ValueOf(ret), nil
	})

	m.AddConversion((**dtos.NextAvailableDockingStation)(nil), (**pb.NextAvailableDockingStation)(nil), func(in reflect.Value) (reflect.Value, error) {

		ret := &pb.NextAvailableDockingStation{}
		v := in.Interface().(*dtos.NextAvailableDockingStation)
		errs := m.Copy(ret, v)
		if len(errs) > 0 {
			return reflect.Zero(in.Type()), errs[0]
		}

		return reflect.ValueOf(ret), nil
	})
}
