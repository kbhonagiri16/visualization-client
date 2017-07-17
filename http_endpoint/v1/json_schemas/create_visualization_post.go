package v1JsonSchema

// VisualizationsCreateJSONSchema describes data expected by app on /visualizations url
const VisualizationsCreateJSONSchema = `{
    "$schema": "http://json-schema.org/schema#",
    "type": "object",
    "properties": {
        "dashboards": {
            "type": "array",
            "items": {
                "type": "object",
				"properties": {
					"name": {
						"type": "string"
					},
					"templateBody": {
						"type": "string"
					},
					"templateName": {
						"type": "string"
					},
					"templateParameters": {
						"type": "object"
					},
					"templateVersion": {
						"type": "integer"
					}
				},
				"required": [
					"name",
					"templateParameters"
				],
				"additionalProperties": false,
				"oneOf": [
					{
						"required": ["templateName", "templateVersion"],
						"not": {"required": ["templateBody"]}
					},
					{
						"allOf": [
							{"not": {"required": ["templateName"]}},
							{"not": {"required": ["templateVersion"]}},
							{"required": ["templateBody"]}
						]
					}
				]
            }
        },
        "name": {
            "type": "string"
        },
        "tags": {
            "type": "object"
        }
    },
    "required": [
        "name",
        "dashboards"
    ],
	"additionalProperties": false
}`
