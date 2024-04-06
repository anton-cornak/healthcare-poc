// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/add": {
            "post": {
                "description": "Add all numbers provided in the payload",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Add numbers",
                "operationId": "add-operation",
                "parameters": [
                    {
                        "description": "Numbers to add",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.AddPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/compute": {
            "post": {
                "description": "Adds all numbers in the 'add' list and subtracts all numbers in the 'subtract' list.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Math Operations"
                ],
                "summary": "Compute result",
                "parameters": [
                    {
                        "description": "Compute payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.ComputePayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/location": {
            "post": {
                "description": "Get the WKT location based on the user's location",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get WKT location",
                "operationId": "location",
                "parameters": [
                    {
                        "description": "User location",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.GetWKTLocationPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/specialist": {
            "post": {
                "description": "Find a specialist based on the user's location, specialty, and radius",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Find specialist",
                "operationId": "find-specialist",
                "parameters": [
                    {
                        "description": "Specialty, radius, and user location",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.FindSpecialistPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/specialties": {
            "post": {
                "description": "Get all specialties",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get specialties",
                "operationId": "specialties",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.GetSpecialtiesResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/subtract": {
            "post": {
                "description": "Subtract all numbers in the 'subtract' list from the 'number'.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Math Operations"
                ],
                "summary": "Subtract numbers",
                "parameters": [
                    {
                        "description": "Numbers to substract from the 'number'",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.SubtractPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.AddPayload": {
            "type": "object",
            "properties": {
                "numbers": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                }
            }
        },
        "handlers.ComputePayload": {
            "type": "object",
            "properties": {
                "add": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "subtract": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                }
            }
        },
        "handlers.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "handlers.FindSpecialistPayload": {
            "type": "object",
            "properties": {
                "radius": {
                    "type": "integer"
                },
                "specialty_id": {
                    "type": "integer"
                },
                "user_location": {
                    "type": "string"
                }
            }
        },
        "handlers.GetSpecialtiesResponse": {
            "type": "object",
            "properties": {
                "specialties": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/types.Specialty"
                    }
                }
            }
        },
        "handlers.GetWKTLocationPayload": {
            "type": "object",
            "properties": {
                "user_location": {
                    "type": "string"
                }
            }
        },
        "handlers.SubtractPayload": {
            "type": "object",
            "properties": {
                "number": {
                    "type": "number"
                },
                "subtract": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                }
            }
        },
        "handlers.SuccessResponse": {
            "type": "object",
            "properties": {
                "result": {
                    "type": "number"
                }
            }
        },
        "types.Specialty": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
