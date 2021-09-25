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
)
