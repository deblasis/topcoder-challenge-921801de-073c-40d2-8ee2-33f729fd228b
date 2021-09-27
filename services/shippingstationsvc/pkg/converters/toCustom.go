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
	pb "deblasis.net/space-traffic-control/gen/proto/go/centralcommandsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/consts"
	m "gopkg.in/jeevatkm/go-model.v1"
)

func RegisterShipResponseToRegisteredShip(src *pb.RegisterShipResponse) *RegisteredShip {
	ret := &RegisteredShip{
		//There are only two known statuses in the domain, since the ship is not docked, we assume that it's `in-flight`. Sane-defaulting but also easy to change if needed right here
		Status: consts.SHIP_STATUS_INFLIGHT,
	}
	m.Copy(ret, src.Ship)

	return ret
}

type RegisteredShip pb.Ship
