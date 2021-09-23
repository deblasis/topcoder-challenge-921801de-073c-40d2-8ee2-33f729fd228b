package e2e_tests

var (
	registerShipResponseSchema = `
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
)
