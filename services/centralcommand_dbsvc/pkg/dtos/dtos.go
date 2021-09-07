package dtos

import (
	"deblasis.net/space-traffic-control/common/errors"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/internal/model"
)

type CreateShipRequest model.Ship
type CreateShipResponse struct {
	Ship *model.Ship `json:"ship"`
	Err  string      `json:"err,omitempty"`
}

type GetAllShipsRequest struct{}
type GetAllShipsResponse struct {
	Ships []*model.Ship `json:"ships"`
	Err   string        `json:"err,omitempty"`
}

type CreateStationRequest model.Station
type CreateStationResponse struct {
	Station *model.Station `json:"station"`
	Err     string         `json:"err,omitempty"`
}

type GetAllStationsRequest struct{}
type GetAllStationsResponse struct {
	Stations []*model.Station `json:"stations"`
	Err      string           `json:"err,omitempty"`
}

// ErrorMessage is for performing the error massage and returning by API
type ErrorMessage struct {
	Error []string `json:"error"`
}

// NewErrorMessage returns ErrorMessage by error string
func NewErrorMessage(err string) ErrorMessage {
	return ErrorMessage{Error: []string{err}}
}

func (r CreateShipResponse) Failed() error    { return errors.Str2err(r.Err) }
func (r GetAllShipsResponse) Failed() error   { return errors.Str2err(r.Err) }
func (r CreateStationResponse) Failed() error { return errors.Str2err(r.Err) }
