package errs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrUnknown    = errors.New("not found")
	ErrBadRequest = errors.New("bad request")

	ErrCannotSelectEntities              = errors.New("cannot select entities")
	ErrCannotSelectEntity                = errors.New("cannot select entity")
	ErrCannotInsertEntity                = errors.New("cannot insert entity")
	ErrCannotInsertAlreadyExistingEntity = errors.New("the entity you are trying to create already exists")

	ErrUnauthenticated = errors.New("cannot identify user")
	ErrUnauthorized    = errors.New("access denied")
)

func Str2err(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func Err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func CompareStatusErrors(error_a, error_b error) bool {
	return status.Convert(error_a).Err() == status.Convert(error_b).Err()
}

func EncodeErrorHTTP(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/vnd.deblasis.spacetrafficcontrol-v1+/json; charset=utf-8")
	GetErrorContainer(ctx).Transport = err

	w.WriteHeader(Err2code(err))

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func ErrorHandlerGRPC(ctx context.Context, err error) {
	GetErrorContainer(ctx).Transport = err
	header := metadata.Pairs(
		"x-http-code", fmt.Sprintf("%v", Err2code(err)),
		"x-stc-error-message", err.Error(),
	)
	grpc.SendHeader(ctx, header)

}

func Err2code(err error) int {
	switch err.Error() {
	case ErrBadRequest.Error(), ErrCannotInsertAlreadyExistingEntity.Error():
		return http.StatusBadRequest
	case ErrUnknown.Error():
		return http.StatusNotFound
	case ErrUnauthenticated.Error(), ErrUnauthorized.Error():
		return http.StatusUnauthorized
	}

	return http.StatusInternalServerError
}

// func ErrorEncoderGRPC(ctx context.Context, err error, grpcResponse *interface{}) (interface{}, error) {
// 	if err != nil {
// 		header := metadata.Pairs("x-http-code", fmt.Sprintf("%v", errs.Err2code(err)))
// 		grpc.SendHeader(ctx, header)
// 	}
// 	return grpcResponse, nil
// }
