//go:build integration
// +build integration

package e2e_tests

import (
	"context"
	"os"

	"github.com/go-kit/kit/log"

	"deblasis.net/space-traffic-control/common/consts"
	centralcommand_dbsvc_v1 "deblasis.net/space-traffic-control/gen/proto/go/centralcommand_dbsvc/v1"
	"deblasis.net/space-traffic-control/services/centralcommand_dbsvc/pkg/transport"
	"github.com/bxcodec/faker/v3"
	"github.com/gavv/httpexpect/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
)

func GetCommandUserToken(client *httpexpect.Expect) string {

	var tokenReq = client.POST(AuthService_Login).WithJSON(loginRequest{
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

	var signupReq signupRequest

	if err := faker.FakeData(&signupReq); err != nil {
		panic(err)
	}
	signupReq.Role = consts.ROLE_STATION

	return commandClient.POST(AuthService_Signup).WithJSON(signupReq).
		Expect().JSON().Path("$.token.token").String().Raw()

}

func GetNewShipUserToken(client *httpexpect.Expect) string {

	var signupReq signupRequest

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

		loginResp = commandClient.POST(AuthService_Login).WithJSON(loginRequest{
			Username: s,
			Password: s,
		}).Expect().JSON()

		if loginResp.Path("$.error").Raw() == nil || loginResp.Path("$.error").Raw() == "" {
			token = loginResp.Path("$.token.token").String().Raw()
			personas[s] = token
		} else {
			personas[s] = commandClient.POST(AuthService_Signup).WithJSON(signupRequest{
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
		loginResp = commandClient.POST(AuthService_Login).WithJSON(loginRequest{
			Username: s,
			Password: s,
		}).Expect().JSON()

		if loginResp.Path("$.error").Raw() == nil || loginResp.Path("$.error").Raw() == "" {
			token = loginResp.Path("$.token.token").String().Raw()
			personas[s] = token
		} else {
			personas[s] = commandClient.POST(AuthService_Signup).WithJSON(signupRequest{
				Username: s,
				Password: s,
				Role:     consts.ROLE_STATION,
			}).Expect().JSON().Path("$.token.token").String().Raw()
		}
	}

}

func CleanupDB(ctx context.Context, logger log.Logger) error {
	var (
		conn *grpc.ClientConn
		err  error
	)
	target := "localhost:9383"
	if envTarget := os.Getenv("CENTRALCOMMAND_AUX_DBENDPOINT"); envTarget != "" {
		target = envTarget
	}

	conn, err = grpc.Dial(target, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return err
	}

	client := transport.NewAuxGrpcClient(conn, logger)
	resp, err := client.Cleanup(ctx, &centralcommand_dbsvc_v1.CleanupRequest{})
	if err != nil {
		return err
	}
	return resp.Error
}
