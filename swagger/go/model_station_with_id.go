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
/*
 * [deblasis] Space Traffic Control
 *
 * The [MCRN](https://expanse.fandom.com/wiki/Martian_Congressional_Republic_Navy) wants to build and deploy new software to all their space stations spread throughout the Solar System. With the exponential increase of trade between the [OPA](https://expanse.fandom.com/wiki/Outer_Planets_Alliance) and the [Earth](https://expanse.fandom.com/wiki/Earth) the legacy systems running on Martian space stations have been having difficulty directing and optimizing traffic.
 *
 * API version: 1.0.0
 * Contact: alex@deblasis.net
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type StationWithId struct {

	Docks []DockWithId `json:"docks"`

	Id string `json:"id"`
}
