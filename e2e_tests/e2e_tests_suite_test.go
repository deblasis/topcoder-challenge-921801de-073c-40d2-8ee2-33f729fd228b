//go:build integration
// +build integration

package e2e_tests

import (
	"context"
	"fmt"
	"os"
	"testing"

	"deblasis.net/space-traffic-control/common/consts"
	. "deblasis.net/space-traffic-control/e2e_tests/utils"
	"github.com/bxcodec/faker/v3"
	"github.com/gavv/httpexpect/v2"
	"github.com/go-kit/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	client *httpexpect.Expect

	personas map[string]string
)

func TestE2ETests(t *testing.T) {
	RegisterFailHandler(Fail)

	target := "http://localhost:8081"
	if envTarget := os.Getenv("APIGATEWAY"); envTarget != "" {
		target = envTarget
	}

	client = httpexpect.New(&ginkgoTestReporter{}, target)
	GinkgoWriter.Write([]byte("\n⏳ Initializing test harness, creating test users and getting their credentials...\n"))
	personas = make(map[string]string)
	personas[Persona_Command_Initial] = GetCommandUserToken(client)

	bootstrapInitialUsers()
	GinkgoWriter.Write([]byte("\n✅ Initialised, running tests\n\n"))

	RunSpecs(t, "E2E Tests Suite")
}

var _ = BeforeSuite(func() {
	Expect(client).NotTo(BeNil())
	Expect(personas[Persona_Command_Initial]).NotTo(BeEmpty())
})

var _ = AfterSuite(func() {

	ctx := context.Background()
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	err := CleanupDB(ctx, logger)
	Expect(err).Should(Not(HaveOccurred()))

})

type ginkgoTestReporter struct{}

func (g ginkgoTestReporter) Errorf(format string, args ...interface{}) {
	Fail(fmt.Sprintf(format, args...))
}

func (g ginkgoTestReporter) Fatalf(format string, args ...interface{}) {
	Fail(fmt.Sprintf(format, args...))
}

func (g ginkgoTestReporter) Logf(format string, args ...interface{}) {
	GinkgoWriter.Write([]byte(fmt.Sprintf(format+"\n", args...)))
}

func GetCommandUserToken(client *httpexpect.Expect) string {

	var tokenReq = client.POST(AuthService_Login).WithJSON(LoginRequest{
		Username: "deblasis", //TODO config?
		Password: "password!",
	}).Expect().JSON().Path("$.token.token")

	Expect(tokenReq.NotNull()).NotTo(BeNil())

	return tokenReq.String().Raw()
}

func GetNewStationUserToken(client *httpexpect.Expect) string {

	commandClient := client.Builder(func(r *httpexpect.Request) {
		r.WithHeader("Authorization", "Bearer "+personas[Persona_Command_Initial])
	})

	var signupReq SignupRequest

	if err := faker.FakeData(&signupReq); err != nil {
		panic(err)
	}
	signupReq.Role = consts.ROLE_STATION

	return commandClient.POST(AuthService_Signup).WithJSON(signupReq).
		Expect().JSON().Path("$.token.token").String().Raw()

}

func GetNewShipUserToken(client *httpexpect.Expect) string {

	var signupReq SignupRequest

	if err := faker.FakeData(&signupReq); err != nil {
		panic(err)
	}
	signupReq.Role = consts.ROLE_SHIP

	return client.POST(AuthService_Signup).WithJSON(signupReq).
		Expect().JSON().Path("$.token.token").String().Raw()

}

func bootstrapInitialUsers() {

	commandClient := client.Builder(func(r *httpexpect.Request) {
		r.WithHeader("Authorization", "Bearer "+personas[Persona_Command_Initial])
	})

	var (
		loginResp *httpexpect.Value
		token     string
	)
	ships := []string{
		Persona_Ship_MilleniumFalcon,
		Persona_Ship_USSEnterprise,
	}
	for _, s := range ships {

		loginResp = commandClient.POST(AuthService_Login).WithJSON(LoginRequest{
			Username: s,
			Password: s,
		}).Expect().JSON()

		if loginResp.Path("$.error").Raw() == nil || loginResp.Path("$.error").Raw() == "" {
			token = loginResp.Path("$.token.token").String().Raw()
			personas[s] = token
		} else {
			personas[s] = commandClient.POST(AuthService_Signup).WithJSON(SignupRequest{
				Username: s,
				Password: s,
				Role:     consts.ROLE_SHIP,
			}).Expect().JSON().Path("$.token.token").String().Raw()
		}
	}

	stations := []string{
		Persona_Station_DeathStar,
		Persona_Station_ISS,
	}
	for _, s := range stations {
		loginResp = commandClient.POST(AuthService_Login).WithJSON(LoginRequest{
			Username: s,
			Password: s,
		}).Expect().JSON()

		if loginResp.Path("$.error").Raw() == nil || loginResp.Path("$.error").Raw() == "" {
			token = loginResp.Path("$.token.token").String().Raw()
			personas[s] = token
		} else {
			personas[s] = commandClient.POST(AuthService_Signup).WithJSON(SignupRequest{
				Username: s,
				Password: s,
				Role:     consts.ROLE_STATION,
			}).Expect().JSON().Path("$.token.token").String().Raw()
		}
	}

}
