// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "description": "API for AI-powered functionality in a language learning app",
    "title": "Tutor",
    "version": "1.0.0"
  },
  "paths": {
    "/chat": {
      "post": {
        "description": "This endpoint receives a user's message and returns the AI's response.",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Sends a message to the AI and receives a response.",
        "parameters": [
          {
            "description": "User message and additional information",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "message": {
                  "description": "The message text sent by the user.",
                  "type": "string"
                },
                "userId": {
                  "description": "Unique identifier for the user.",
                  "type": "string"
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "type": "object",
              "properties": {
                "reply": {
                  "description": "AI's response to the user's message.",
                  "type": "string"
                }
              }
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/error"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          },
          "default": {
            "description": "error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "error": {
      "type": "object",
      "properties": {
        "code": {
          "type": "integer",
          "format": "int64"
        },
        "message": {
          "type": "string"
        }
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "description": "API for AI-powered functionality in a language learning app",
    "title": "Tutor",
    "version": "1.0.0"
  },
  "paths": {
    "/chat": {
      "post": {
        "description": "This endpoint receives a user's message and returns the AI's response.",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Sends a message to the AI and receives a response.",
        "parameters": [
          {
            "description": "User message and additional information",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "message": {
                  "description": "The message text sent by the user.",
                  "type": "string"
                },
                "userId": {
                  "description": "Unique identifier for the user.",
                  "type": "string"
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "type": "object",
              "properties": {
                "reply": {
                  "description": "AI's response to the user's message.",
                  "type": "string"
                }
              }
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/error"
            }
          },
          "500": {
            "description": "Internal server error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          },
          "default": {
            "description": "error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "error": {
      "type": "object",
      "properties": {
        "code": {
          "type": "integer",
          "format": "int64"
        },
        "message": {
          "type": "string"
        }
      }
    }
  }
}`))
}
