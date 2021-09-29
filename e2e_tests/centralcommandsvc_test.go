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
	"encoding/json"
	"fmt"
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

			shipClients    map[string]*httpexpect.Expect
			stationClients map[string]*httpexpect.Expect
			commandClient  *httpexpect.Expect
		)
		BeforeEach(func() {
			shipClients = make(map[string]*httpexpect.Expect)
			stationClients = make(map[string]*httpexpect.Expect)
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
				shipClients[Persona_Ship_USSEnterprise] = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+personas[Persona_Ship_USSEnterprise])
				})
				if err := faker.FakeData(&registerStationReq); err != nil {
					panic(err)
				}
			})
			It("should fail returning 401", func() {
				shipClients[Persona_Ship_USSEnterprise].POST(CentralCommandService_RegisterStation).
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

				stationClients["new"] = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+newStationToken)
				})
				if err := faker.FakeData(&registerStationReq); err != nil {
					panic(err)
				}
			})
			It("should succeed returning 200 and fail on subsequent attempts with 401", func() {
				validCall := stationClients["new"].POST(CentralCommandService_RegisterStation).
					WithJSON(registerStationReq).Expect()

				validCall.Status(http.StatusOK)

				validCall.JSON().Schema(RegisterStationResponseSchema)

				stationClients["new"].POST(CentralCommandService_RegisterStation).
					WithJSON(registerStationReq).Expect().Status(http.StatusBadRequest)
				stationClients["new"].POST(CentralCommandService_RegisterStation).
					WithJSON(registerStationReq).Expect().Status(http.StatusBadRequest)
				stationClients["new"].POST(CentralCommandService_RegisterStation).
					WithJSON(registerStationReq).Expect().Status(http.StatusBadRequest)

			})
		})
	})
	Describe("Register Ship", func() {
		var (
			registerShipReq RegisterShipRequest

			shipClients   map[string]*httpexpect.Expect
			commandClient *httpexpect.Expect
		)
		BeforeEach(func() {
			shipClients = make(map[string]*httpexpect.Expect)
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
				shipClients[Persona_Station_DeathStar] = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+personas[Persona_Station_DeathStar])
				})
				if err := faker.FakeData(&registerShipReq); err != nil {
					panic(err)
				}
			})
			It("should fail returning 401", func() {
				shipClients[Persona_Station_DeathStar].POST(CentralCommandService_RegisterShip).
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

				shipClients["new"] = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+newShipToken)
				})
				if err := faker.FakeData(&registerShipReq); err != nil {
					panic(err)
				}
			})
			It("should succeed returning 200, empty body and fail on subsequent attempts with 401", func() {
				validCall := shipClients["new"].POST(CentralCommandService_RegisterShip).
					WithJSON(registerShipReq).Expect()

				validCall.Status(http.StatusOK)

				validCall.Body().Empty()
				Expect(validCall.Raw().ContentLength).To(BeEquivalentTo(0))

				shipClients["new"].POST(CentralCommandService_RegisterShip).
					WithJSON(registerShipReq).Expect().Status(http.StatusBadRequest)
				shipClients["new"].POST(CentralCommandService_RegisterShip).
					WithJSON(registerShipReq).Expect().Status(http.StatusBadRequest)
				shipClients["new"].POST(CentralCommandService_RegisterShip).
					WithJSON(registerShipReq).Expect().Status(http.StatusBadRequest)

			})
		})
	})

	Describe("List Stations", func() {

		var (
			shipClients    map[string]*httpexpect.Expect
			stationClients map[string]*httpexpect.Expect
			commandClient  *httpexpect.Expect
		)

		BeforeEach(func() {
			shipClients = make(map[string]*httpexpect.Expect)
			stationClients = make(map[string]*httpexpect.Expect)
		})

		When("a token is not provided", func() {

			It("should fail returning 401", func() {
				client.GET(CentralCommandService_AllStations).
					Expect().Status(http.StatusUnauthorized)
			})
		})
		When("a Station token is provided", func() {
			BeforeEach(func() {
				stationClients[Persona_Station_DeathStar] = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+personas[Persona_Station_DeathStar])
				})
			})
			It("should fail returning 400", func() {
				stationClients[Persona_Station_DeathStar].GET(CentralCommandService_AllStations).
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
				shipClients["new"] = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+newShipToken)
				})
			})
			It("should succeed returning 200", func() {
				validCall := shipClients["new"].GET(CentralCommandService_AllStations).
					Expect()

				validCall.Status(http.StatusOK)

				validCall.JSON().Schema(GetAllStationsResponseSchema)

			})

			When("there are registered stations", func() {
				var (
					n_registered_stations   = 5
					stationTokens           = make([]string, n_registered_stations)
					getStationReqsByUserId  = make(map[string]interface{}, n_registered_stations)
					getStationRespsByUserId = make(map[string]interface{}, n_registered_stations)
				)

				BeforeEach(func() {
					stationClients = make(map[string]*httpexpect.Expect)
					CleanupDB(ctx, logger)

					for i := 0; i < n_registered_stations; i++ {
						stationTokens[i] = GetNewStationUserToken(client)

						clientClaims, _ := jwtHandler.ExtractClaims(stationTokens[i])
						userId := clientClaims.UserId

						stationClients[fmt.Sprintf("station_%v", i)] = client.Builder(func(r *httpexpect.Request) {
							r.WithHeader("Authorization", "Bearer "+stationTokens[i])
						})
						getStationReqsByUserId[userId] = &RegisterStationRequest{}
						faker.FakeData(getStationReqsByUserId[userId])

						resp := stationClients[fmt.Sprintf("station_%v", i)].POST(CentralCommandService_RegisterStation).
							WithJSON(getStationReqsByUserId[userId]).Expect()

						getStationRespsByUserId[userId] = resp.Body().Raw()
					}
				})

				When("a command token is supplied", func() {

					It("should return all the stations correctly", func() {

						validCall := commandClient.GET(CentralCommandService_AllStations).
							Expect()

						validCall.Status(http.StatusOK)
						validCall.JSON().Schema(GetAllStationsResponseSchema)

						validCall.JSON().Array().Length().Equal(len(getStationReqsByUserId))

						//validCall.JSON().Array().Empty()

						expected := make(map[string]*RegisterStationResponse, 0)
						for userId, respJson := range getStationRespsByUserId {

							resp := &RegisterStationResponse{}
							json.Unmarshal([]byte(respJson.(string)), resp)

							station := &RegisterStationResponse{
								Id:           resp.Id,
								Capacity:     getStationReqsByUserId[userId].(*RegisterStationRequest).Capacity,
								UsedCapacity: resp.UsedCapacity,
								Docks:        make([]*DockResp, 0),
							}
							for _, dock := range resp.Docks {
								station.Docks = append(station.Docks, &DockResp{
									Id:              dock.Id,
									NumDockingPorts: dock.NumDockingPorts,
									Occupied:        dock.Occupied,
									Weight:          dock.Weight,
								})
							}
							expected[userId] = station
						}

						totalFound := 0
						for _, val := range validCall.JSON().Array().Iter() {
							station := val.Raw().(map[string]interface{})
							e := expected[station["id"].(string)]
							if e.Id == station["id"].(string) && BarelyEqual(e.Capacity, station["capacity"].(float64)) && BarelyEqual(e.UsedCapacity, station["usedCapacity"].(float64)) {
								docksCount := 0
								for _, dock := range station["docks"].([]interface{}) {
									for _, edock := range e.Docks {
										if edock.Id == dock.(map[string]interface{})["id"].(string) && BarelyEqual(float64(edock.NumDockingPorts), dock.(map[string]interface{})["numDockingPorts"].(float64)) {
											docksCount++
											break
										}
									}
								}
								if docksCount == len(e.Docks) {
									totalFound++
								}
							}
						}
						Expect(totalFound).To(Equal(len(getStationReqsByUserId)))
					})
				})

				When("a registered ship token is supplied", func() {
					BeforeEach(func() {
						stationClients = make(map[string]*httpexpect.Expect)
						CleanupDB(ctx, logger)

						stations := []struct {
							capacity float64
							docks    []*Dock
						}{
							{
								capacity: 10,
								docks: []*Dock{
									{
										NumDockingPorts: 1,
									},
									{
										NumDockingPorts: 9,
									},
								},
							},
							{
								capacity: 4,
								docks: []*Dock{
									{
										NumDockingPorts: 2,
									},
									{
										NumDockingPorts: 2,
									},
								},
							},
							{
								capacity: 1,
								docks: []*Dock{
									{
										NumDockingPorts: 2,
									},
								},
							},
						}
						for i, tC := range stations {
							tC := tC

							stationTokens[i] = GetNewStationUserToken(client)

							clientClaims, _ := jwtHandler.ExtractClaims(stationTokens[i])
							userId := clientClaims.UserId

							stationClients[fmt.Sprintf("station_%v", i)] = client.Builder(func(r *httpexpect.Request) {
								r.WithHeader("Authorization", "Bearer "+stationTokens[i])
							})
							getStationReqsByUserId[userId] = &RegisterStationRequest{
								Capacity: tC.capacity,
								Docks:    tC.docks,
							}

							stationClients[fmt.Sprintf("station_%v", i)].POST(CentralCommandService_RegisterStation).
								WithJSON(getStationReqsByUserId[userId]).Expect()
						}
					})

					It("should return all stations when the ship it's so small it would fit in any available station", func() {
						newShipToken := GetNewShipUserToken(client)
						shipClients["small"] = client.Builder(func(r *httpexpect.Request) {
							r.WithHeader("Authorization", "Bearer "+newShipToken)
						})
						shipClients["small"].POST(CentralCommandService_RegisterShip).
							WithJSON(&RegisterShipRequest{
								Weight: 1,
							}).Expect()

						validCall := shipClients["small"].GET(CentralCommandService_AllStations).
							Expect()

						validCall.Status(http.StatusOK)
						validCall.JSON().Schema(GetAllStationsResponseSchema)

						validCall.JSON().Array().Length().Equal(3)

					})

					It("should return only one station when the ship is too big to fit in all stations", func() {
						newShipToken := GetNewShipUserToken(client)
						shipClients["medium"] = client.Builder(func(r *httpexpect.Request) {
							r.WithHeader("Authorization", "Bearer "+newShipToken)
						})
						shipClients["medium"].POST(CentralCommandService_RegisterShip).
							WithJSON(&RegisterShipRequest{
								Weight: 5,
							}).Expect()
						validCall := shipClients["medium"].GET(CentralCommandService_AllStations).
							Expect()

						validCall.Status(http.StatusOK)
						validCall.JSON().Schema(GetAllStationsResponseSchema)

						validCall.JSON().Array().Length().Equal(1)

					})
					It("should return no stations since the ship is to big to fit anywhere", func() {
						newShipToken := GetNewShipUserToken(client)
						shipClients["huge"] = client.Builder(func(r *httpexpect.Request) {
							r.WithHeader("Authorization", "Bearer "+newShipToken)
						})
						shipClients["huge"].POST(CentralCommandService_RegisterShip).
							WithJSON(&RegisterShipRequest{
								Weight: 10000000,
							}).Expect()
						validCall := shipClients["huge"].GET(CentralCommandService_AllStations).
							Expect()

						validCall.Status(http.StatusOK)
						validCall.JSON().Schema(GetAllStationsResponseSchema)

						validCall.JSON().Array().Length().Equal(0)

					})
				})

			})
		})

	})
	Describe("List Ships", func() {
		var (
			shipClients    map[string]*httpexpect.Expect
			stationClients map[string]*httpexpect.Expect
			commandClient  *httpexpect.Expect
		)

		BeforeEach(func() {
			shipClients = make(map[string]*httpexpect.Expect)
			stationClients = make(map[string]*httpexpect.Expect)
		})

		When("a token is not provided", func() {

			It("should fail returning 401", func() {
				client.GET(CentralCommandService_AllShips).
					Expect().Status(http.StatusUnauthorized)
			})
		})
		When("a Station token is provided", func() {
			BeforeEach(func() {
				stationClients[Persona_Station_DeathStar] = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+personas[Persona_Station_DeathStar])
				})
			})
			It("should fail returning 401", func() {
				stationClients[Persona_Station_DeathStar].GET(CentralCommandService_AllShips).
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

			When("there are registered ships", func() {
				var (
					n_registered_ships = 5
					shipTokens         = make([]string, n_registered_ships)
					regShipReqs        = make([]interface{}, n_registered_ships)
				)

				BeforeEach(func() {
					shipClients = make(map[string]*httpexpect.Expect)
					CleanupDB(ctx, logger)

					for i := 0; i < n_registered_ships; i++ {
						shipTokens[i] = GetNewShipUserToken(client)
						shipClients[fmt.Sprintf("ship_%v", i)] = client.Builder(func(r *httpexpect.Request) {
							r.WithHeader("Authorization", "Bearer "+shipTokens[i])
						})
						regShipReqs[i] = &RegisterShipRequest{}
						faker.FakeData(regShipReqs[i])

						sr := shipClients[fmt.Sprintf("ship_%v", i)].POST(CentralCommandService_RegisterShip).
							WithJSON(regShipReqs[i]).Expect()
						sr.Body().Empty()
					}
				})

				It("should return them correctly", func() {

					validCall := commandClient.GET(CentralCommandService_AllShips).
						Expect()

					validCall.Status(http.StatusOK)
					validCall.JSON().Schema(GetAllShipsResponseSchema)

					validCall.JSON().Array().Length().Equal(len(regShipReqs))

					expected := make([]*RegisterShipResponse, 0)
					for i, token := range shipTokens {
						clientClaims, _ := jwtHandler.ExtractClaims(token)
						expected = append(expected, &RegisterShipResponse{
							Id:     clientClaims.UserId,
							Weight: regShipReqs[i].(*RegisterShipRequest).Weight,
							Status: "in-flight",
						})
					}
					//this doesn't work because of floats
					//validCall.JSON().Array().Contains(expected)
					totalFound := 0
					for _, e := range expected {
						for _, val := range validCall.JSON().Array().Iter() {
							ship := val.Raw().(map[string]interface{})
							if e.Id == ship["id"].(string) && BarelyEqual(float64(e.Weight), ship["weight"].(float64)) {
								totalFound++
								break
							}
						}
					}
					Expect(totalFound).To(Equal(len(regShipReqs)))
				})
			})
		})
		When("a Ship token is provided", func() {
			BeforeEach(func() {
				shipClients[Persona_Ship_USSEnterprise] = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+personas[Persona_Ship_USSEnterprise])
				})
			})
			It("should fail returning 401", func() {
				shipClients[Persona_Ship_USSEnterprise].GET(CentralCommandService_AllShips).
					Expect().Status(http.StatusUnauthorized)
			})
		})
	})

})
