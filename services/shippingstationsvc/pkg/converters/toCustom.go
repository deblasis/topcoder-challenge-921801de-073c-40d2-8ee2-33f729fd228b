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
