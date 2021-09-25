//go:build integration
// +build integration

package e2e_tests

const (
	AuthService_Signup                    = "/user/signup"
	AuthService_Login                     = "/auth/login"
	CentralCommandService_RegisterStation = "/centcom/station/register"
	CentralCommandService_AllStations     = "/centcom/station/all"
	CentralCommandService_RegisterShip    = "/centcom/ship/register"
	CentralCommandService_AllShips        = "/centcom/ship/all"
	ShippingStationService_RequestLanding = "/shipping-station/request-landing"
	ShippingStationService_Land           = "/shipping-station/land"
)
