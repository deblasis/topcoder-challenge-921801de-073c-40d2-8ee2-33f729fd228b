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
package errs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	common_v1 "deblasis.net/space-traffic-control/gen/proto/go/v1"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	m "gopkg.in/jeevatkm/go-model.v1"
)

var (
	ErrException  = errors.New("exception")
	ErrUnknown    = errors.New("not found")
	ErrBadRequest = errors.New("bad request")

	ErrValidationFailed = errors.New("validation failed")

	ErrCannotSelectEntities              = errors.New("cannot select entities")
	ErrCannotSelectEntity                = errors.New("cannot select entity")
	ErrCannotInsertEntity                = errors.New("cannot insert entity")
	ErrCannotInsertAlreadyExistingEntity = errors.New("the entity you are trying to create already exists")

	ErrUnauthenticated = errors.New("cannot identify user")
	ErrUnauthorized    = errors.New("access denied")
)

func EncodeErrorHTTP(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if e, ok := err.(*Err); ok {
		w.WriteHeader(int(e.Code))
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": e.Message,
		})
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func TranportErrorHandler(logger log.Logger) func(ctx context.Context, err error) {

	return func(ctx context.Context, err error) {
		logger.Log("err", err)
		if e, ok := err.(*Err); ok {
			GetErrorContainer(ctx).Transport = e
			header := metadata.Pairs(
				"x-http-code", fmt.Sprintf("%v", e.Code),
			)
			grpc.SendHeader(ctx, header)
		}
	}
}

// func Err2code(err error) int {

// 	switch {
// 	case Is(err, ErrBadRequest) ||
// 		Is(err, ErrCannotInsertAlreadyExistingEntity) ||
// 		Is(err, ErrValidationFailed):
// 		return http.StatusBadRequest
// 	case Is(err, ErrUnknown):
// 		return http.StatusNotFound
// 	case Is(err, ErrUnauthenticated) ||
// 		Is(err, ErrUnauthorized):
// 		return http.StatusUnauthorized
// 	}

// 	return http.StatusInternalServerError
// }

func IsNil(e error) bool {
	var err *Err
	if errors.As(e, &err) {
		return err == NilErr
	}
	var perr *common_v1.Error
	if errors.As(e, &perr) {
		return perr == common_v1.NilErr
	}
	return e == nil
}

var NilErr = (*Err)(nil)

type Err struct {
	Code    int32  `json:"-"`
	Message string `json:"message,omitempty"`
	Err     error  `json:"-,omitempty" model:"-,omitempty"`
}

func NewError(code int32, message string, err error) error {
	return &Err{
		Message: message,
		Code:    code,
		Err:     err,
	}
}

func (e *Err) Error() string {
	if e == (*Err)(nil) {
		return ""
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *Err) Unwrap() error {
	return e.Err
}

// Returns the inner most Error
func (err *Err) Dig() *Err {
	var ew *Err
	if errors.As(err.Err, &ew) {
		// Recursively digs until wrapper error is not in which case it will stop
		return ew.Dig()
	}
	return err
}

func MarshalJSON(err error) []byte {
	val, ok := err.(*Err)
	if !ok {
		return MarshalJSON(err)
	}
	j, err := json.Marshal(val)
	if err != nil {
		return MarshalJSON(&Err{
			Code: val.Code,
			Message: fmt.Sprint("can't convert error to JSON Base64:",
				base64.StdEncoding.EncodeToString([]byte(val.Error()))),
		})
	}
	return j
}

func (e Err) ToProtoV1() *common_v1.Error {
	if e == (Err{}) {
		return nil
	}
	return &common_v1.Error{
		Code:    e.Code,
		Message: e.Message,
	}
}

func ToProtoV1(e error) *common_v1.Error {
	if e == nil {
		return nil
	}

	var ee *Err
	if errors.As(e, &ee) {
		return ee.ToProtoV1()
	}

	ret := &common_v1.Error{}
	ret.Code = http.StatusInternalServerError
	ret.Message = e.Error()

	return ret
}

func FromProtoV1(e *common_v1.Error) *Err {
	if e == (*common_v1.Error)(nil) {
		return nil
	}
	ret := &Err{}
	m.Copy(ret, e)

	return ret
}

func InjectGrpcStatusCode(ctx context.Context, resp interface{}) {
	if failer, ok := resp.(endpoint.Failer); ok {
		e := failer.Failed()
		InjectGrpcErrorStatusCode(ctx, e, http.StatusInternalServerError)
	}
}

func InjectGrpcErrorStatusCode(ctx context.Context, e error, fallbackCode int) {

	var code int
	var err *Err
	var perr *common_v1.Error

	if !IsNil(e) {
		if errors.As(e, &err) {
			code = int(err.Code)
		}
		if errors.As(e, &perr) {
			code = int(perr.Code)
		}
		if err, ok := status.FromError(e); ok {
			code = runtime.HTTPStatusFromCode(err.Code())
		}
		if code == 0 {
			code = fallbackCode
		}
		header := metadata.Pairs(
			"x-http-code", fmt.Sprintf("%v", code),
		)
		grpc.SendHeader(ctx, header)
	}

}
