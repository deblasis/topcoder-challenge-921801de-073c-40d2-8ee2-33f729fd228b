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
//go:build integration
// +build integration

package e2e_tests

import (
	"net/http"
	"time"

	. "deblasis.net/space-traffic-control/e2e_tests/utils"
	"github.com/bxcodec/faker/v3"
	"github.com/gavv/httpexpect/v2"
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
