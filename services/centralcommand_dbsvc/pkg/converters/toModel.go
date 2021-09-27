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
