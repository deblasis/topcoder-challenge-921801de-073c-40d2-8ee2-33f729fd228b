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
