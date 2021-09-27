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
