//
// Copyright 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package converters

// import (
// 	"strings"

// 	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
// )

// func ProtoToDTORole(src pb.SignupRequest_Role) string {
// 	return strings.Title(strings.TrimLeft(strings.ToLower(src.String()), "role_"))
// }

import (
	pb "deblasis.net/space-traffic-control/gen/proto/go/authsvc/v1"
	"deblasis.net/space-traffic-control/services/auth_dbsvc/pkg/dtos"
	m "gopkg.in/jeevatkm/go-model.v1"
)

func SignupRequestToDBDto(src *pb.SignupRequest) *dtos.CreateUserRequest {
	ret := &dtos.CreateUserRequest{}
	if errs := m.Copy(ret, src); len(errs) > 0 {
		//panic(errs[0])
	}
	return ret
}
