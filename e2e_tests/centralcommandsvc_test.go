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

	. "deblasis.net/space-traffic-control/e2e_tests/utils"
	"github.com/bxcodec/faker/v3"
	"github.com/gavv/httpexpect/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CentralCommandSvc", func() {

	Describe("httpClient", func() {
		It("should be initialized successfully", func() {
			Expect(client).NotTo(BeNil())
		})
	})

	Describe("Register Station", func() {
		var (
			registerStationReq RegisterStationRequest

			shipClient    *httpexpect.Expect
			stationClient *httpexpect.Expect
			commandClient *httpexpect.Expect
		)
		BeforeEach(func() {
			registerStationReq = RegisterStationRequest{}
		})

		When("a token is not provided", func() {
			BeforeEach(func() {
				if err := faker.FakeData(&registerStationReq); err != nil {
					panic(err)
				}
			})
			It("should fail returning 401", func() {
				client.POST(CentralCommandService_RegisterStation).
					WithJSON(registerStationReq).Expect().Status(http.StatusUnauthorized)
			})
		})
		When("a Ship token is provided", func() {
			BeforeEach(func() {
				shipClient = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+personas[Persona_Ship_USSEnterprise])
				})
				if err := faker.FakeData(&registerStationReq); err != nil {
					panic(err)
				}
			})
			It("should fail returning 401", func() {
				shipClient.POST(CentralCommandService_RegisterStation).
					WithJSON(registerStationReq).Expect().Status(http.StatusUnauthorized)
			})
		})
		When("a Command token is provided", func() {
			BeforeEach(func() {
				commandClient = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+personas[Persona_Command_Initial])
				})
				if err := faker.FakeData(&registerStationReq); err != nil {
					panic(err)
				}
			})
			It("should fail returning 401", func() {
				commandClient.POST(CentralCommandService_RegisterStation).
					WithJSON(registerStationReq).Expect().Status(http.StatusUnauthorized)
			})
		})
		When("a Station token is provided", func() {
			BeforeEach(func() {
				newStationToken := GetNewStationUserToken(client)

				stationClient = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+newStationToken)
				})
				if err := faker.FakeData(&registerStationReq); err != nil {
					panic(err)
				}
			})
			It("should succeed returning 200 and fail on subsequent attempts with 401", func() {
				validCall := stationClient.POST(CentralCommandService_RegisterStation).
					WithJSON(registerStationReq).Expect()

				validCall.Status(http.StatusOK)

				validCall.JSON().Schema(RegisterStationResponseSchema)

				stationClient.POST(CentralCommandService_RegisterStation).
					WithJSON(registerStationReq).Expect().Status(http.StatusBadRequest)
				stationClient.POST(CentralCommandService_RegisterStation).
					WithJSON(registerStationReq).Expect().Status(http.StatusBadRequest)
				stationClient.POST(CentralCommandService_RegisterStation).
					WithJSON(registerStationReq).Expect().Status(http.StatusBadRequest)

			})
		})
	})
	Describe("Register Ship", func() {
		var (
			registerShipReq RegisterShipRequest

			shipClient    *httpexpect.Expect
			stationClient *httpexpect.Expect
			commandClient *httpexpect.Expect
		)
		BeforeEach(func() {
			registerShipReq = RegisterShipRequest{}
		})

		When("a token is not provided", func() {
			BeforeEach(func() {
				if err := faker.FakeData(&registerShipReq); err != nil {
					panic(err)
				}
			})
			It("should fail returning 401", func() {
				client.POST(CentralCommandService_RegisterShip).
					WithJSON(registerShipReq).Expect().Status(http.StatusUnauthorized)
			})
		})
		When("a Station token is provided", func() {
			BeforeEach(func() {
				shipClient = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+personas[Persona_Station_DeathStar])
				})
				if err := faker.FakeData(&registerShipReq); err != nil {
					panic(err)
				}
			})
			It("should fail returning 401", func() {
				shipClient.POST(CentralCommandService_RegisterShip).
					WithJSON(registerShipReq).Expect().Status(http.StatusUnauthorized)
			})
		})
		When("a Command token is provided", func() {
			BeforeEach(func() {
				commandClient = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+personas[Persona_Command_Initial])
				})
				if err := faker.FakeData(&registerShipReq); err != nil {
					panic(err)
				}
			})
			It("should fail returning 401", func() {
				commandClient.POST(CentralCommandService_RegisterShip).
					WithJSON(registerShipReq).Expect().Status(http.StatusUnauthorized)
			})
		})
		When("a Ship token is provided", func() {
			BeforeEach(func() {
				newShipToken := GetNewShipUserToken(client)

				stationClient = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+newShipToken)
				})
				if err := faker.FakeData(&registerShipReq); err != nil {
					panic(err)
				}
			})
			It("should succeed returning 200, empty body and fail on subsequent attempts with 401", func() {
				validCall := stationClient.POST(CentralCommandService_RegisterShip).
					WithJSON(registerShipReq).Expect()

				validCall.Status(http.StatusOK)

				validCall.Body().Empty()
				Expect(validCall.Raw().ContentLength).To(BeEquivalentTo(0))

				stationClient.POST(CentralCommandService_RegisterShip).
					WithJSON(registerShipReq).Expect().Status(http.StatusBadRequest)
				stationClient.POST(CentralCommandService_RegisterShip).
					WithJSON(registerShipReq).Expect().Status(http.StatusBadRequest)
				stationClient.POST(CentralCommandService_RegisterShip).
					WithJSON(registerShipReq).Expect().Status(http.StatusBadRequest)

			})
		})
	})

	Describe("List Stations", func() {

		var (
			shipClient    *httpexpect.Expect
			stationClient *httpexpect.Expect
			commandClient *httpexpect.Expect
		)

		When("a token is not provided", func() {

			It("should fail returning 401", func() {
				client.GET(CentralCommandService_AllStations).
					Expect().Status(http.StatusUnauthorized)
			})
		})
		When("a Station token is provided", func() {
			BeforeEach(func() {
				stationClient = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+personas[Persona_Station_DeathStar])
				})
			})
			It("should fail returning 400", func() {
				stationClient.GET(CentralCommandService_AllStations).
					Expect().Status(http.StatusBadRequest)
			})
		})
		When("a Command token is provided", func() {
			BeforeEach(func() {
				commandClient = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+personas[Persona_Command_Initial])
				})
			})
			It("should succeed returning 200", func() {
				validCall := commandClient.GET(CentralCommandService_AllStations).
					Expect()

				validCall.Status(http.StatusOK)
				validCall.JSON().Schema(GetAllStationsResponseSchema)

			})
		})
		When("a Ship token is provided", func() {
			BeforeEach(func() {
				newShipToken := GetNewShipUserToken(client)
				shipClient = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+newShipToken)
				})
			})
			It("should succeed returning 200", func() {
				validCall := shipClient.GET(CentralCommandService_AllStations).
					Expect()

				validCall.Status(http.StatusOK)

				validCall.JSON().Schema(GetAllStationsResponseSchema)

			})
		})

	})
	Describe("List Ships", func() {
		var (
			shipClient    *httpexpect.Expect
			stationClient *httpexpect.Expect
			commandClient *httpexpect.Expect
		)

		When("a token is not provided", func() {

			It("should fail returning 401", func() {
				client.GET(CentralCommandService_AllShips).
					Expect().Status(http.StatusUnauthorized)
			})
		})
		When("a Station token is provided", func() {
			BeforeEach(func() {
				stationClient = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+personas[Persona_Station_DeathStar])
				})
			})
			It("should fail returning 401", func() {
				stationClient.GET(CentralCommandService_AllShips).
					Expect().Status(http.StatusUnauthorized)
			})
		})
		When("a Command token is provided", func() {
			BeforeEach(func() {
				commandClient = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+personas[Persona_Command_Initial])
				})
			})
			It("should succeed returning 200", func() {
				validCall := commandClient.GET(CentralCommandService_AllShips).
					Expect()

				validCall.Status(http.StatusOK)
				validCall.JSON().Schema(GetAllShipsResponseSchema)

			})
		})
		When("a Ship token is provided", func() {
			BeforeEach(func() {
				shipClient = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+personas[Persona_Ship_USSEnterprise])
				})
			})
			It("should fail returning 401", func() {
				shipClient.GET(CentralCommandService_AllShips).
					Expect().Status(http.StatusUnauthorized)
			})
		})
	})

})
