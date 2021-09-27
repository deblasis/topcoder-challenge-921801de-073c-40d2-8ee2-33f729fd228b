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
