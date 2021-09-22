package e2e_tests_test

import (
	"fmt"
	"testing"

	"deblasis.net/space-traffic-control/common/consts"
	"github.com/gavv/httpexpect/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	client *httpexpect.Expect

	personas map[string]string
)

func TestE2ETests(t *testing.T) {
	RegisterFailHandler(Fail)

	client = httpexpect.New(&ginkgoTestReporter{}, "http://localhost:8081")
	personas = make(map[string]string)
	personas[Persona_Command_Initial] = GetCommandUserToken(client)

	bootstrapInitialUsers()
	RunSpecs(t, "E2E Tests Suite")
}

var _ = BeforeSuite(func() {
	Expect(client).NotTo(BeNil())
	Expect(personas[Persona_Command_Initial]).NotTo(BeEmpty())
})

func GetCommandUserToken(client *httpexpect.Expect) string {

	var tokenReq = client.POST("/auth/login").WithJSON(loginRequest{
		Username: "deblasis",
		Password: "password!",
	}).Expect().JSON().Path("$.token.token")

	Expect(tokenReq.NotNull()).NotTo(BeNil())

	return tokenReq.String().Raw()
}

func GetStationUserToken(client *httpexpect.Expect) string {

	commandClient := client.Builder(func(r *httpexpect.Request) {
		r.WithHeader("Authorization", "Bearer "+personas[Persona_Command_Initial])
	})

	return commandClient.POST("/user/signup").WithJSON(signupRequest{
		Username: "deblasis",
		Password: "password!",
	}).Expect().JSON().Path("$.token.token").String().Raw()

}

type signupRequest struct {
	Username string `json:"username,omitempty" faker:"username"`
	Password string `json:"password,omitempty" faker:"password"`
	Role     string `json:"role,omitempty" faker:"-"` //faker:"oneof: Ship, Station, Command"
}
type loginRequest struct {
	Username string `json:"username,omitempty" faker:"username"`
	Password string `json:"password,omitempty" faker:"password"`
}

type ginkgoTestReporter struct{}

func (g ginkgoTestReporter) Errorf(format string, args ...interface{}) {
	Fail(fmt.Sprintf(format, args...))
}

func (g ginkgoTestReporter) Fatalf(format string, args ...interface{}) {
	Fail(fmt.Sprintf(format, args...))
}

func (g ginkgoTestReporter) Logf(format string, args ...interface{}) {
	GinkgoWriter.Write([]byte(fmt.Sprintf(format, args...)))
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

		loginResp = commandClient.POST("/auth/login").WithJSON(loginRequest{
			Username: s,
			Password: s,
		}).Expect().JSON()

		if loginResp.Path("$.error").Raw() == nil || loginResp.Path("$.error").Raw() == "" {
			token = loginResp.Path("$.token.token").String().Raw()
			personas[s] = token
		} else {
			personas[s] = commandClient.POST("/user/signup").WithJSON(signupRequest{
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
		loginResp = commandClient.POST("/auth/login").WithJSON(loginRequest{
			Username: s,
			Password: s,
		}).Expect().JSON()

		if loginResp.Path("$.error").Raw() == nil || loginResp.Path("$.error").Raw() == "" {
			token = loginResp.Path("$.token.token").String().Raw()
			personas[s] = token
		} else {
			personas[s] = commandClient.POST("/user/signup").WithJSON(signupRequest{
				Username: s,
				Password: s,
				Role:     consts.ROLE_STATION,
			}).Expect().JSON().Path("$.token.token").String().Raw()
		}
	}

}

const (
	Persona_Command_Initial = "Persona_Command_Initial"

	Persona_Ship_MilleniumFalcon = "Persona_Ship_MilleniumFalcon"
	Persona_Ship_USSEnterprise   = "Persona_Ship_USSEnterprise"

	Persona_Station_DeathStar = "Persona_Station_DeathStar"
	Persona_Station_ISS       = "Persona_Station_ISS"
)
