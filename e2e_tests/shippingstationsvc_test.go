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

package e2e_tests

import (
	"context"
	"net/http"
	"os"
	"time"

	. "deblasis.net/space-traffic-control/e2e_tests/utils"
	"github.com/bxcodec/faker/v3"
	"github.com/gavv/httpexpect/v2"
	"github.com/go-kit/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ShippingStationSvc", func() {

	Describe("httpClient", func() {
		It("should be initialized successfully", func() {
			Expect(client).NotTo(BeNil())
		})
	})

	Describe("Request Landing", func() {
		var (
			requestLandingReq RequestLandingRequest

			shipClients    map[string]*httpexpect.Expect
			stationClients map[string]*httpexpect.Expect
			commandClient  *httpexpect.Expect
		)
		BeforeEach(func() {
			requestLandingReq = RequestLandingRequest{}
			shipClients = make(map[string]*httpexpect.Expect)
			stationClients = make(map[string]*httpexpect.Expect)

			commandClient = client.Builder(func(r *httpexpect.Request) {
				r.WithHeader("Authorization", "Bearer "+personas[Persona_Command_Initial])
			})

			shipClients[Persona_Ship_USSEnterprise] = client.Builder(func(r *httpexpect.Request) {
				r.WithHeader("Authorization", "Bearer "+personas[Persona_Ship_USSEnterprise])
			})
			shipClients[Persona_Ship_MilleniumFalcon] = client.Builder(func(r *httpexpect.Request) {
				r.WithHeader("Authorization", "Bearer "+personas[Persona_Ship_MilleniumFalcon])
			})

			stationClients[Persona_Station_ISS] = client.Builder(func(r *httpexpect.Request) {
				r.WithHeader("Authorization", "Bearer "+personas[Persona_Station_ISS])
			})
			stationClients[Persona_Station_DeathStar] = client.Builder(func(r *httpexpect.Request) {
				r.WithHeader("Authorization", "Bearer "+personas[Persona_Station_DeathStar])
			})

		})

		Context("there are no registrations of any kind yet", func() {

			When("a token is not provided", func() {
				BeforeEach(func() {
					if err := faker.FakeData(&requestLandingReq); err != nil {
						panic(err)
					}
				})
				It("should fail returning 401", func() {
					client.POST(ShippingStationService_RequestLanding).
						WithJSON(requestLandingReq).Expect().Status(http.StatusUnauthorized)
				})
			})

			When("a Ship token is provided", func() {
				BeforeEach(func() {
					if err := faker.FakeData(&requestLandingReq); err != nil {
						panic(err)
					}
				})
				//This could have been a 404 but 503 seems more appropriate since it's not the user that's in the wrong but the server
				It("should fail returning 503", func() {
					shipClients[Persona_Ship_USSEnterprise].POST(ShippingStationService_RequestLanding).
						WithJSON(requestLandingReq).Expect().Status(http.StatusServiceUnavailable)
				})
			})
			When("a Station token is provided", func() {
				BeforeEach(func() {
					if err := faker.FakeData(&requestLandingReq); err != nil {
						panic(err)
					}
				})
				It("should fail returning 401", func() {
					stationClients[Persona_Station_ISS].POST(ShippingStationService_RequestLanding).
						WithJSON(requestLandingReq).Expect().Status(http.StatusUnauthorized)
				})
			})
			When("a Command token is provided", func() {
				BeforeEach(func() {
					if err := faker.FakeData(&requestLandingReq); err != nil {
						panic(err)
					}
				})
				It("should fail returning 401", func() {
					commandClient.POST(ShippingStationService_RequestLanding).
						WithJSON(requestLandingReq).Expect().Status(http.StatusUnauthorized)
				})
			})
		})

		When("there is a station registered", func() {
			BeforeEach(func() {
				ctx := context.Background()
				logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
				CleanupDB(ctx, logger)

				newStationToken := GetNewStationUserToken(client)

				registerStationReq := &RegisterStationRequest{
					Capacity: 15.50,
					Docks: []*Dock{
						{NumDockingPorts: 5},
					},
				}
				stationClients[Persona_Station_ISS] = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+newStationToken)
				})

				stationClients[Persona_Station_ISS].POST(CentralCommandService_RegisterStation).
					WithJSON(registerStationReq).Expect().Status(http.StatusOK)
			})
			It("should fail returning 503 because the ship is not registered", func() {
				requestLandingReq = RequestLandingRequest{Time: 10}
				shipClients[Persona_Ship_USSEnterprise].POST(ShippingStationService_RequestLanding).
					WithJSON(requestLandingReq).Expect().Status(http.StatusServiceUnavailable)
			})
		})

		When("there is a station and a ship registered", func() {
			When("there is capacity", func() {
				BeforeEach(func() {
					ctx := context.Background()
					logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
					CleanupDB(ctx, logger)

					newStationToken := GetNewStationUserToken(client)

					registerStationReq := &RegisterStationRequest{
						Capacity: 15.50,
						Docks: []*Dock{
							{NumDockingPorts: 5},
						},
					}
					stationClients[Persona_Station_ISS] = client.Builder(func(r *httpexpect.Request) {
						r.WithHeader("Authorization", "Bearer "+newStationToken)
					})

					stationClients[Persona_Station_ISS].POST(CentralCommandService_RegisterStation).
						WithJSON(registerStationReq).Expect().Status(http.StatusOK)

					registerShipReq := &RegisterShipRequest{Weight: 10}
					shipClients[Persona_Ship_USSEnterprise].POST(CentralCommandService_RegisterShip).
						WithJSON(registerShipReq).Expect().Status(http.StatusOK)

				})
				It("should succeed with a `land` command", func() {
					requestLandingReq = RequestLandingRequest{Time: 10}
					shipClients[Persona_Ship_USSEnterprise].POST(ShippingStationService_RequestLanding).
						WithJSON(requestLandingReq).Expect().
						JSON().Schema(RequestLandingLandCommandResponseSchema)
				})
			})
			When("there isn't enough capacity", func() {
				BeforeEach(func() {
					ctx := context.Background()
					logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
					CleanupDB(ctx, logger)

					registerStationReq := &RegisterStationRequest{
						Capacity: 15.50,
						Docks: []*Dock{
							{NumDockingPorts: 5},
						},
					}

					stationClients[Persona_Station_ISS].POST(CentralCommandService_RegisterStation).
						WithJSON(registerStationReq).Expect().Status(http.StatusOK)

					registerShipReq := &RegisterShipRequest{
						Weight: 20,
					}
					shipClients[Persona_Ship_USSEnterprise].POST(CentralCommandService_RegisterShip).
						WithJSON(registerShipReq).Expect().Status(http.StatusOK)

				})
				It("should succeed with a `wait` command", func() {
					requestLandingReq = RequestLandingRequest{Time: 10}
					shipClients[Persona_Ship_USSEnterprise].POST(ShippingStationService_RequestLanding).
						WithJSON(requestLandingReq).Expect().
						JSON().Schema(RequestLandingWaitCommandResponseSchema)
				})
			})
		})

		When("there is capacity but there are no available docks", func() {
			BeforeEach(func() {
				ctx := context.Background()
				logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
				CleanupDB(ctx, logger)

				stationClients[Persona_Station_ISS].POST(CentralCommandService_RegisterStation).
					WithJSON(&RegisterStationRequest{
						Capacity: 15.50,
						Docks: []*Dock{
							{NumDockingPorts: 1},
						},
					}).Expect().Status(http.StatusOK)

				shipClients[Persona_Ship_USSEnterprise].POST(CentralCommandService_RegisterShip).
					WithJSON(&RegisterShipRequest{
						Weight: 10,
					}).Expect().Status(http.StatusOK)

				requestLandingResponse := shipClients[Persona_Ship_USSEnterprise].POST(ShippingStationService_RequestLanding).
					WithJSON(RequestLandingRequest{Time: 120}).Expect().
					JSON()

				requestLandingResponse.Schema(RequestLandingLandCommandResponseSchema)
				dockIdForLanding := requestLandingResponse.Path("$.dockingStation").String().Raw()

				shipClients[Persona_Ship_USSEnterprise].POST(ShippingStationService_Land).
					WithJSON(LandRequest{Time: 5, DockId: dockIdForLanding}).Expect().Status(http.StatusOK)

				shipClients[Persona_Ship_MilleniumFalcon].POST(CentralCommandService_RegisterShip).
					WithJSON(&RegisterShipRequest{
						Weight: 5,
					}).Expect().Status(http.StatusOK)

				time.Sleep(100 * time.Millisecond)
				shipClients[Persona_Ship_MilleniumFalcon].POST(ShippingStationService_RequestLanding).
					WithJSON(RequestLandingRequest{Time: 120}).Expect().
					JSON().Object().ContainsMap(map[string]interface{}{
					"command": "wait",
				})

			})
			It("should succeed with a `wait` command", func() {
				shipClients[Persona_Ship_MilleniumFalcon].POST(ShippingStationService_RequestLanding).
					WithJSON(RequestLandingRequest{Time: 10}).Expect().
					JSON().Schema(RequestLandingWaitCommandResponseSchema).
					Object().ContainsMap(map[string]interface{}{
					"command":  "wait",
					"duration": 5,
				})
			})
			It("should return a shorter wait on subsequent calls", func() {
				time.Sleep(2 * time.Second)
				shipClients[Persona_Ship_MilleniumFalcon].POST(ShippingStationService_RequestLanding).
					WithJSON(RequestLandingRequest{Time: 10}).Expect().
					JSON().Schema(RequestLandingWaitCommandResponseSchema).
					Object().ContainsMap(map[string]interface{}{
					"command":  "wait",
					"duration": 3,
				})
			})
			It("should return a `land` command after the previous ship left", func() {
				time.Sleep(6 * time.Second)
				shipClients[Persona_Ship_MilleniumFalcon].POST(ShippingStationService_RequestLanding).
					WithJSON(RequestLandingRequest{Time: 10}).Expect().
					JSON().Schema(RequestLandingLandCommandResponseSchema)
			})
		})
	})

})
