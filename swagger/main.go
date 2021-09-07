/*
 * [deblasis] Space Traffic Control
 *
 * The [MCRN](https://expanse.fandom.com/wiki/Martian_Congressional_Republic_Navy) wants to build and deploy new software to all their space stations spread throughout the Solar System. With the exponential increase of trade between the [OPA](https://expanse.fandom.com/wiki/Outer_Planets_Alliance) and the [Earth](https://expanse.fandom.com/wiki/Earth) the legacy systems running on Martian space stations have been having difficulty directing and optimizing traffic.
 *
 * API version: 1.0.0
 * Contact: alex@deblasis.net
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package main

import (
	"log"
	"net/http"

	// WARNING!
	// Change this to a fully-qualified import path
	// once you place this file into your project.
	// For example,
	//
	//    sw "github.com/myname/myrepo/go"
	//
	sw "deblasis.net/space-traffic-control/swagger/go"
)

func main() {
	log.Printf("Server started")

	router := sw.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}