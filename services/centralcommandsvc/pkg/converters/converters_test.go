package converters_test

import (
	"fmt"

	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/dtos"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/converters"
	"deblasis.net/space-traffic-control/services/centralcommandsvc/pkg/dtos"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	m "gopkg.in/jeevatkm/go-model.v1"
)

var _ = Describe("Converters", func() {

	BeforeSuite(func() {

		converters.ConfigMappers()

	})

	Describe("Custom converters", func() {
		var (
			testCases []struct {
				desc         string
				src          interface{}
				dest         interface{}
				expectedDest interface{}
			}
		)

		Describe("Struct to Struct", func() {

			//src dtos.GetAllShipsResponse) *dtos.GetAllShipsResponse
			testCases = []struct {
				desc         string
				src          interface{}
				dest         interface{}
				expectedDest interface{}
			}{
				{
					desc: "Ship",
					src: &dtos.Ship{
						Id:     "test",
						Status: "test",
						Weight: 10,
					},
					dest: &dtos.Ship{},
					expectedDest: &dtos.Ship{
						Id:     "test",
						Status: "",
						Weight: 10,
					},
				},
				{
					desc: "Dock",
					src: &dtos.Dock{
						Id:              "test",
						StationId:       "",
						NumDockingPorts: 0,
						Occupied:        0,
						Weight:          10,
					},
					dest: &dtos.Dock{},
					expectedDest: &dtos.Dock{
						Id:              "test",
						StationId:       "",
						NumDockingPorts: 0,
						Occupied:        0,
						Weight:          10,
					},
				},
				{
					desc: "GetAllShipsResponse",
					src: &dtos.GetAllShipsResponse{
						Ships: []dtos.Ship{
							{
								Id:     "test",
								Status: "",
								Weight: 10,
							},
						},
						Err: "",
					},
					dest: &dtos.GetAllShipsResponse{
						Ships: []dtos.Ship{
							{
								Id:     "test",
								Status: "",
								Weight: 10,
							},
						},
						Err: "",
					},
					expectedDest: nil,
				},
			}
			for _, tC := range testCases {
				tC := tC

				When(fmt.Sprintf("It's a %v", tC.desc), func() {

					It("Should match", func() {
						//act
						errs := m.Copy(tC.dest, tC.src)
						//defer GinkgoRecover()
						//assert
						Expect(errs).To(BeNil())
						Expect(tC.dest).To(BeEquivalentTo(tC.expectedDest))
					})
				})

			}
		})

	})

})
