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
        "/songs": {
            "get": {
                "description": "Get all songs with filtering and pagination",
                "produces": [
                    "application/json"
                ],
                "summary": "Get all songs",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Song"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Add a new song",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Add a new song",
                "parameters": [
                    {
                        "description": "Song data",
                        "name": "song",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/services.NewSongRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Song"
                        }
                    }
                }
            }
        },
        "/songs/{id}": {
            "get": {
                "description": "Get a song by ID",
                "produces": [
                    "application/json"
                ],
                "summary": "Get a song by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Song ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Song"
                        }
                    }
                }
            },
            "put": {
                "description": "Update an existing song",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Update a song",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Song ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Song data",
                        "name": "song",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Song"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Song"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a song by ID",
                "produces": [
                    "application/json"
                ],
                "summary": "Delete a song",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Song ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Song": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "type": "string"
                },
                "group": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "link": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "release_date": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "services.NewSongRequest": {
            "type": "object",
            "required": [
                "group",
                "song"
            ],
            "properties": {
                "group": {
                    "type": "string"
                },
                "song": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Music Library API",
	Description:      "API для управления библиотекой песен.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
