package e2e_tests_test

import (
	"fmt"
	"net/http"

	"deblasis.net/space-traffic-control/common/consts"
	"github.com/bxcodec/faker/v3"
	"github.com/gavv/httpexpect/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AuthSvc", func() {

	It("Testing harness should be initialized successfully", func() {
		Expect(client).NotTo(BeNil())
	})

	Describe("Signup", func() {

		var (
			signupReq signupRequest
		)

		BeforeEach(func() {
			signupReq = signupRequest{}
		})
		Context("unauthenticated user", func() {
			When("trying to sign up a Station user", func() {
				BeforeEach(func() {
					if err := faker.FakeData(&signupReq); err != nil {
						panic(err)
					}
					signupReq.Role = consts.ROLE_STATION
				})
				It("should be unauthorized", func() {
					client.POST("/user/signup").WithJSON(signupReq).Expect().Status(http.StatusUnauthorized)
				})
			})
			When("trying to sign up a Command user", func() {
				BeforeEach(func() {
					if err := faker.FakeData(&signupReq); err != nil {
						panic(err)
					}
					signupReq.Role = consts.ROLE_COMMAND
				})
				It("should be unauthorized", func() {
					client.POST("/user/signup").WithJSON(signupReq).Expect().Status(http.StatusUnauthorized)
				})
			})
			When("trying to signing up with malformed role", func() {
				BeforeEach(func() {
					if err := faker.FakeData(&signupReq); err != nil {
						panic(err)
					}
					signupReq.Role = "Pirate"
				})
				It("should fail", func() {
					client.POST("/user/signup").WithJSON(signupReq).Expect().Status(http.StatusUnauthorized)
				})
			})

			When("trying to signing up a Ship user", func() {
				BeforeEach(func() {
					if err := faker.FakeData(&signupReq); err != nil {
						panic(err)
					}
					signupReq.Role = consts.ROLE_SHIP
				})
				It("should succeed", func() {
					client.POST("/user/signup").WithJSON(signupReq).Expect().Status(http.StatusOK)
				})
			})

			When("trying to signing up a Ship user with empty credentials", func() {
				BeforeEach(func() {
					signupReq.Role = consts.ROLE_SHIP
				})
				It("should fail with validation error", func() {
					client.POST("/user/signup").WithJSON(signupReq).Expect().Status(http.StatusBadRequest)
				})
			})

			When("trying to signing up a Ship user twice", func() {
				BeforeEach(func() {
					if err := faker.FakeData(&signupReq); err != nil {
						panic(err)
					}
					signupReq.Role = consts.ROLE_SHIP
				})
				It("should succeed at first and fail on subsequent requests", func() {
					client.POST("/user/signup").WithJSON(signupReq).Expect().Status(http.StatusOK)
					client.POST("/user/signup").WithJSON(signupReq).Expect().Status(http.StatusBadRequest)
					client.POST("/user/signup").WithJSON(signupReq).Expect().Status(http.StatusBadRequest)
					client.POST("/user/signup").WithJSON(signupReq).Expect().Status(http.StatusBadRequest)
				})
			})

		})

		Context("logged in as Command user", func() {

			var commandClient *httpexpect.Expect

			BeforeEach(func() {
				commandClient = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+personas[Persona_Command_Initial])
				})
			})

			testCases := []struct {
				role           string
				expectedstatus int
			}{
				{
					role:           consts.ROLE_SHIP,
					expectedstatus: http.StatusOK,
				},
				{
					role:           consts.ROLE_STATION,
					expectedstatus: http.StatusOK,
				},
				{
					role:           consts.ROLE_COMMAND,
					expectedstatus: http.StatusOK,
				},
			}
			for _, tC := range testCases {
				When(fmt.Sprintf("trying to sign up a %v user", tC.role), func() {
					BeforeEach(func() {
						if err := faker.FakeData(&signupReq); err != nil {
							panic(err)
						}
						signupReq.Role = consts.ROLE_SHIP
					})
					It("should succeed", func() {
						commandClient.POST("/user/signup").
							WithJSON(signupReq).Expect().Status(tC.expectedstatus)
					})
				})
			}

		})

	})

	// Context("logged in as Station user", func() {

	// 	var stationClient *httpexpect.Expect

	// 	BeforeEach(func() {
	// 		stationClient = client.Builder(func(r *httpexpect.Request) {
	// 			r.WithHeader("Authorization", "Bearer "+ )
	// 		})
	// 	})

	// 	testCases := []struct {
	// 		role           string
	// 		expectedstatus int
	// 	}{
	// 		{
	// 			role:           consts.ROLE_SHIP,
	// 			expectedstatus: http.StatusOK,
	// 		},
	// 		{
	// 			role:           consts.ROLE_STATION,
	// 			expectedstatus: http.StatusOK,
	// 		},
	// 		{
	// 			role:           consts.ROLE_COMMAND,
	// 			expectedstatus: http.StatusOK,
	// 		},
	// 	}
	// 	for _, tC := range testCases {
	// 		When(fmt.Sprintf("trying to sign up a %v user", tC.role), func() {
	// 			BeforeEach(func() {
	// 				if err := faker.FakeData(&signupReq); err != nil {
	// 					panic(err)
	// 				}
	// 				signupReq.Role = consts.ROLE_SHIP
	// 			})
	// 			It("should succeed", func() {
	// 				stationClient.POST("/user/signup").
	// 					WithJSON(signupReq).Expect().Status(tC.expectedstatus)
	// 			})
	// 		})
	// 	}

	// })

	// })

})
