/*
 * [deblasis] Space Traffic Control
 *
 * The [MCRN](https://expanse.fandom.com/wiki/Martian_Congressional_Republic_Navy) wants to build and deploy new software to all their space stations spread throughout the Solar System. With the exponential increase of trade between the [OPA](https://expanse.fandom.com/wiki/Outer_Planets_Alliance) and the [Earth](https://expanse.fandom.com/wiki/Earth) the legacy systems running on Martian space stations have been having difficulty directing and optimizing traffic.
 *
 * API version: 1.0.0
 * Contact: alex@deblasis.net
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package v1

import (
	"net/http"
)

func ShipLand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ShipRequestLanding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
