//go:build integration
// +build integration

package e2e_tests

import (
	"fmt"
	"os"
	"testing"

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
