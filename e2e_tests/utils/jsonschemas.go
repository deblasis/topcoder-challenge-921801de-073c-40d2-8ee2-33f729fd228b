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
//go:build integration
// +build integration

package utils

var (
	RegisterStationResponseSchema = `
	{
    "$schema": "http://json-schema.org/draft-07/schema",
    "$id": "http://example.com/example.json",
    "type": "object",
    "required": [
        "id",
        "docks"
    ],
    "properties": {
        "id": {
            "$id": "#/properties/id",
            "type": "string"
        },
        "docks": {
            "$id": "#/properties/docks",
            "type": "array",
            "additionalItems": true,
            "items": {
                "$id": "#/properties/docks/items",
                "anyOf": [
                    {
                        "$id": "#/properties/docks/items/anyOf/0",
                        "type": "object",
                        "required": [
                            "id",
                            "numDockingPorts"
                        ],
                        "properties": {
                            "id": {
                                "$id": "#/properties/docks/items/anyOf/0/properties/id",
                                "type": "string"
                            },
                            "numDockingPorts": {
                                "$id": "#/properties/docks/items/anyOf/0/properties/numDockingPorts",
                                "type": "integer"
                            }
                        },
                        "additionalProperties": false
                    }
                ]
            }
        }
    },
    "additionalProperties": false
}
`

	GetAllStationsResponseSchema = `
{
  "$schema": "http://json-schema.org/draft-07/schema",
  "$id": "http://example.com/example.json",
  "type": "array",
  "default": [],
  "additionalItems": true,
  "items": {
      "$id": "#/items",
      "anyOf": [
          {
              "$id": "#/items/anyOf/0",
              "type": "object",
              "default": {},
              "required": [
                  "id",
                  "capacity",
                  "usedCapacity",
                  "docks"
              ],
              "additionalProperties": false,
              "properties": {
                  "id": {
                      "$id": "#/items/anyOf/0/properties/id",
                      "type": "string",
                      "default": ""
                  },
                  "capacity": {
                      "$id": "#/items/anyOf/0/properties/capacity",
                      "type": "number",
                      "default": 0.0
                  },
                  "usedCapacity": {
                      "$id": "#/items/anyOf/0/properties/usedCapacity",
                      "type": "number",
                      "default": 0.0
                  },
                  "docks": {
                      "$id": "#/items/anyOf/0/properties/docks",
                      "type": "array",
                      "default": [],
                      "additionalItems": true,
                      "items": {
                          "$id": "#/items/anyOf/0/properties/docks/items",
                          "anyOf": [
                              {
                                  "$id": "#/items/anyOf/0/properties/docks/items/anyOf/0",
                                  "type": "object",
                                  "default": {},
                                  "required": [
                                      "id",
                                      "numDockingPorts",
                                      "weight"
                                  ],
                                  "additionalProperties": false,
                                  "properties": {
                                      "id": {
                                          "$id": "#/items/anyOf/0/properties/docks/items/anyOf/0/properties/id",
                                          "type": "string",
                                          "default": ""
                                      },
                                      "numDockingPorts": {
                                          "$id": "#/items/anyOf/0/properties/docks/items/anyOf/0/properties/numDockingPorts",
                                          "type": "integer",
                                          "default": 0
                                      },
                                      "occupied": {
                                        "$id": "#/items/anyOf/0/properties/docks/items/anyOf/0/properties/occupied",
                                        "type": "integer",
                                        "default": 0
                                      },
                                      "weight": {
                                          "$id": "#/items/anyOf/0/properties/docks/items/anyOf/0/properties/weight",
                                          "type": "number",
                                          "default": 0.0
                                      }
                                  }
                              }
                          ]
                      }
                  }
              }
          }
      ]
  }
}
`

	GetAllShipsResponseSchema = `
  {
    "$schema": "http://json-schema.org/draft-07/schema",
    "$id": "http://example.com/example.json",
    "type": "array",
    "additionalItems": true,
    "items": {
        "$id": "#/items",
        "anyOf": [
            {
                "$id": "#/items/anyOf/0",
                "type": "object",
                "title": "The first anyOf schema",
                "required": [
                    "id",
                    "status",
                    "weight"
                ],
                "properties": {
                    "id": {
                        "$id": "#/items/anyOf/0/properties/id",
                        "type": "string"
                    },
                    "status": {
                        "$id": "#/items/anyOf/0/properties/status",
                        "type": "string"
                    },
                    "weight": {
                        "$id": "#/items/anyOf/0/properties/weight",
                        "type": "number"
                    }
                },
                "additionalProperties": false
            }
        ]
    }
}
`

	RequestLandingLandCommandResponseSchema = `
{
  "$schema": "http://json-schema.org/draft-07/schema",
  "$id": "http://example.com/example.json",
  "type": "object",
  "required": [
      "command",
      "dockingStation"
  ],
  "properties": {
      "command": {
          "$id": "#/properties/command",
          "type": "string",
          "enum": [
              "land"
          ]
      },
      "dockingStation": {
          "$id": "#/properties/dockingStation",
          "type": "string"
      }
  },
  "additionalProperties": false
}`

	RequestLandingWaitCommandResponseSchema = `
{
  "$schema": "http://json-schema.org/draft-07/schema",
  "$id": "http://example.com/example.json",
  "type": "object",
  "required": [
      "command",
      "duration"
  ],
  "properties": {
      "command": {
          "$id": "#/properties/command",
          "type": "string",
          "enum": [
              "wait"
          ]
      },
      "duration": {
          "$id": "#/properties/duration",
          "type": "number"
      }
  },
  "additionalProperties": true
}`
)
