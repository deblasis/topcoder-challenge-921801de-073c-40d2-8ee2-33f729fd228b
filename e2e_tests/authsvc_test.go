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
	"fmt"
	"net/http"

	"deblasis.net/space-traffic-control/common/consts"
	. "deblasis.net/space-traffic-control/e2e_tests/utils"
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
			signupReq SignupRequest
		)

		BeforeEach(func() {
			signupReq = SignupRequest{}
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
					client.POST(AuthService_Signup).WithJSON(signupReq).Expect().Status(http.StatusUnauthorized)
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
					client.POST(AuthService_Signup).WithJSON(signupReq).Expect().Status(http.StatusUnauthorized)
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
					client.POST(AuthService_Signup).WithJSON(signupReq).Expect().Status(http.StatusUnauthorized)
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
					client.POST(AuthService_Signup).WithJSON(signupReq).Expect().Status(http.StatusOK)
				})
			})

			When("trying to signing up a Ship user with empty credentials", func() {
				BeforeEach(func() {
					signupReq.Role = consts.ROLE_SHIP
				})
				It("should fail with validation error", func() {
					client.POST(AuthService_Signup).WithJSON(signupReq).Expect().Status(http.StatusBadRequest)
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
					client.POST(AuthService_Signup).WithJSON(signupReq).Expect().Status(http.StatusOK)
					client.POST(AuthService_Signup).WithJSON(signupReq).Expect().Status(http.StatusBadRequest)
					client.POST(AuthService_Signup).WithJSON(signupReq).Expect().Status(http.StatusBadRequest)
					client.POST(AuthService_Signup).WithJSON(signupReq).Expect().Status(http.StatusBadRequest)
				})
			})

		})

		Context("logged in as Ship user", func() {

			var shipClient *httpexpect.Expect

			BeforeEach(func() {
				shipClient = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+personas[Persona_Ship_USSEnterprise])
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
					expectedstatus: http.StatusUnauthorized,
				},
				{
					role:           consts.ROLE_COMMAND,
					expectedstatus: http.StatusUnauthorized,
				},
			}
			for _, tC := range testCases {
				tC := tC
				When(fmt.Sprintf("trying to sign up a %v user", tC.role), func() {
					BeforeEach(func() {
						if err := faker.FakeData(&signupReq); err != nil {
							panic(err)
						}
						signupReq.Role = tC.role
					})
					It(fmt.Sprintf("should return status %v", tC.expectedstatus), func() {
						shipClient.POST(AuthService_Signup).
							WithJSON(signupReq).Expect().Status(tC.expectedstatus)
					})
				})
			}

		})

		Context("logged in as Station user", func() {

			var stationClient *httpexpect.Expect

			BeforeEach(func() {
				stationClient = client.Builder(func(r *httpexpect.Request) {
					r.WithHeader("Authorization", "Bearer "+personas[Persona_Station_ISS])
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
					expectedstatus: http.StatusUnauthorized,
				},
				{
					role:           consts.ROLE_COMMAND,
					expectedstatus: http.StatusUnauthorized,
				},
			}
			for _, tC := range testCases {
				tC := tC
				When(fmt.Sprintf("trying to sign up a %v user", tC.role), func() {
					BeforeEach(func() {
						if err := faker.FakeData(&signupReq); err != nil {
							panic(err)
						}
						signupReq.Role = tC.role
					})
					It(fmt.Sprintf("should return status %v", tC.expectedstatus), func() {
						stationClient.POST(AuthService_Signup).
							WithJSON(signupReq).Expect().Status(tC.expectedstatus)
					})
				})
			}

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
				tC := tC
				When(fmt.Sprintf("trying to sign up a %v user", tC.role), func() {
					BeforeEach(func() {
						if err := faker.FakeData(&signupReq); err != nil {
							panic(err)
						}
						signupReq.Role = tC.role
					})
					It(fmt.Sprintf("should return status %v", tC.expectedstatus), func() {
						commandClient.POST(AuthService_Signup).
							WithJSON(signupReq).Expect().Status(tC.expectedstatus)
					})
				})
			}

		})

	})

	Describe("Login", func() {
		var (
			loginReq LoginRequest
		)

		BeforeEach(func() {
			loginReq = LoginRequest{}
		})
		When("providing empty credentials", func() {
			It("should fail with 400", func() {
				client.POST(AuthService_Login).WithJSON(loginReq).Expect().Status(http.StatusBadRequest)
			})
		})
		When("providing wrong credentials", func() {
			BeforeEach(func() {
				if err := faker.FakeData(&loginReq); err != nil {
					panic(err)
				}
			})
			It("should fail with 401", func() {
				client.POST(AuthService_Login).WithJSON(loginReq).Expect().Status(http.StatusUnauthorized)
			})
		})

		When("providing correct credentials", func() {
			BeforeEach(func() {
				loginReq = LoginRequest{
					Username: Persona_Ship_MilleniumFalcon,
					Password: Persona_Ship_MilleniumFalcon,
				}
			})
			It("should succeed with 200", func() {
				client.POST(AuthService_Login).WithJSON(loginReq).Expect().Status(http.StatusOK)
			})
			It("should return token", func() {
				client.POST(AuthService_Login).WithJSON(loginReq).Expect().JSON().Path("$.token").NotNull()
			})
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
	// 				stationClient.POST(AuthService_Signup).
	// 					WithJSON(signupReq).Expect().Status(tC.expectedstatus)
	// 			})
	// 		})
	// 	}

	// })

	// })

})
