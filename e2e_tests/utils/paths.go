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
//go:build integration
// +build integration

package utils

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
