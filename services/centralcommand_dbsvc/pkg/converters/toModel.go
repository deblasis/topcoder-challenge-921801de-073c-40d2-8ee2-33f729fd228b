package converters

import (
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	m "gopkg.in/jeevatkm/go-model.v1"
)

func ShipToModel(src dtos.Ship) model.Ship {
	return model.Ship{
		Id:     src.Id,
		Status: src.Status,
		Weight: src.Weight,
	}
}

func StationToModel(src dtos.Station) model.Station {
	ret := &model.Station{}
	m.Copy(ret, src)
	return *ret
}

func DockToModel(src dtos.Dock) model.Dock {
	return model.Dock{
		Id:              src.Id,
		StationId:       src.StationId,
		NumDockingPorts: src.NumDockingPorts,
		Occupied:        src.Occupied,
		Weight:          src.Weight,
	}
}

func DocksToModel(src []dtos.Dock) []model.Dock {
	var ret []model.Dock
	for _, x := range src {
		ret = append(ret, DockToModel(x))
	}
	return ret
}
