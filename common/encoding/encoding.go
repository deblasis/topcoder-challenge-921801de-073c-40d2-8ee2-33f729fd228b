package encoding

import (
	"context"
	"encoding/json"
	"net/http"

	"deblasis.net/space-traffic-control/common/errors"
)

func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/vnd.deblasis.spacetrafficcontrol-v1+/json; charset=utf-8")
	switch err {
	case errors.ErrUnknown:
		w.WriteHeader(http.StatusNotFound)
	case errors.ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
