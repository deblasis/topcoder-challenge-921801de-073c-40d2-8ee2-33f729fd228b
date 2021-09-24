//go:build integration
// +build integration

package e2e_tests

import (
	"net/http"

	"github.com/bxcodec/faker/v3"
	"github.com/gavv/httpexpect/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CentralCommandSvc", func() {

	It("Testing harness should be initialized successfully", func() {
		Expect(client).NotTo(BeNil())
	})

	Describe("Register Station", func() {
		var (
			registerStationReq registerStationRequest

			shipClient    *httpexpect.Expect
			stationClient *httpexpect.Expect
			commandClient *httpexpect.Expect
		)
		BeforeEach(func() {
			registerStationReq = registerStationRequest{}
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

				validCall.JSON().Schema(registerShipResponseSchema)

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
			registerShipReq registerShipRequest

			shipClient    *httpexpect.Expect
			stationClient *httpexpect.Expect
			commandClient *httpexpect.Expect
		)
		BeforeEach(func() {
			registerShipReq = registerShipRequest{}
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
				shipClient = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+personas[Persona_Station_DeathStar])
				})
			})
			It("should fail returning 400", func() {
				shipClient.GET(CentralCommandService_AllStations).
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
				validCall.JSON().Schema(getAllStationsResponseSchema)

			})
		})
		When("a Ship token is provided", func() {
			BeforeEach(func() {
				newShipToken := GetNewShipUserToken(client)
				stationClient = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+newShipToken)
				})
			})
			It("should succeed returning 200", func() {
				validCall := stationClient.GET(CentralCommandService_AllStations).
					Expect()

				validCall.Status(http.StatusOK)

				validCall.JSON().Schema(getAllStationsResponseSchema)

				// validCall.Body().Empty()
				// Expect(validCall.Raw().ContentLength).To(BeEquivalentTo(0))

				// stationClient.GET(CentralCommandService_AllStations).
				// 	Expect().Status(http.StatusBadRequest)
				// stationClient.GET(CentralCommandService_AllStations).
				// 	Expect().Status(http.StatusBadRequest)
				// stationClient.GET(CentralCommandService_AllStations).
				// 	Expect().Status(http.StatusBadRequest)

			})
		})

	})
	Describe("List Ships", func() {

	})

})