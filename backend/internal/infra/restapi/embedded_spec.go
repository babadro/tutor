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
    "/chat_messages": {
      "post": {
        "description": "This endpoint receives a user's message and returns the AI's response.",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Sends a message to the AI and receives a response.",
        "operationId": "SendChatMessage",
        "parameters": [
          {
            "description": "User message and additional information",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "chatId",
                "text",
                "timestamp"
              ],
              "properties": {
                "chatId": {
                  "description": "The chat ID.",
                  "type": "string"
                },
                "text": {
                  "description": "The message text sent by the user.",
                  "type": "string"
                },
                "timestamp": {
                  "description": "The timestamp of the message.",
                  "type": "integer",
                  "format": "int64"
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
                },
                "timestamp": {
                  "type": "integer",
                  "format": "int64"
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
          "401": {
            "description": "unauthorized",
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
    },
    "/chat_messages/{chatId}": {
      "get": {
        "description": "Get chat messages",
        "produces": [
          "application/json"
        ],
        "operationId": "GetChatMessages",
        "parameters": [
          {
            "type": "string",
            "name": "chatId",
            "in": "path",
            "required": true
          },
          {
            "type": "integer",
            "format": "int32",
            "default": 10,
            "name": "limit",
            "in": "query"
          },
          {
            "type": "integer",
            "format": "int64",
            "default": 0,
            "name": "before_timestamp",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "A list of chat messages",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/ChatMessage"
              }
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/error"
            }
          },
          "401": {
            "description": "unauthorized",
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
    },
    "/voice_messages": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Processes a voice message and returns a response.",
        "operationId": "SendVoiceMessage",
        "parameters": [
          {
            "description": "User message containing a voice message url",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "voiceMessageUrl": {
                  "description": "URL of the voice message file",
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
                "voiceMessageTranscript": {
                  "description": "Transcription of the input voice message.",
                  "type": "string"
                },
                "voiceMessageUrl": {
                  "description": "URL of the voice message input file.",
                  "type": "string"
                },
                "voiceResponseTranscript": {
                  "description": "Text transcription of the voice response.",
                  "type": "string"
                },
                "voiceResponseUrl": {
                  "description": "URL to the voice response file.",
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
          "401": {
            "description": "unauthorized",
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
    "ChatMessage": {
      "type": "object",
      "properties": {
        "curUsr": {
          "type": "boolean",
          "x-go-name": "IsFromCurrentUser"
        },
        "text": {
          "type": "string"
        },
        "timestamp": {
          "type": "integer",
          "format": "int64"
        },
        "userId": {
          "type": "string"
        }
      }
    },
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
  },
  "securityDefinitions": {
    "key": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header",
      "authorizationUrl": ""
    }
  },
  "security": [
    {
      "key": []
    }
  ]
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "description": "API for AI-powered functionality in a language learning app",
    "title": "Tutor",
    "version": "1.0.0"
  },
  "paths": {
    "/chat_messages": {
      "post": {
        "description": "This endpoint receives a user's message and returns the AI's response.",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Sends a message to the AI and receives a response.",
        "operationId": "SendChatMessage",
        "parameters": [
          {
            "description": "User message and additional information",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "chatId",
                "text",
                "timestamp"
              ],
              "properties": {
                "chatId": {
                  "description": "The chat ID.",
                  "type": "string"
                },
                "text": {
                  "description": "The message text sent by the user.",
                  "type": "string"
                },
                "timestamp": {
                  "description": "The timestamp of the message.",
                  "type": "integer",
                  "format": "int64"
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
                },
                "timestamp": {
                  "type": "integer",
                  "format": "int64"
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
          "401": {
            "description": "unauthorized",
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
    },
    "/chat_messages/{chatId}": {
      "get": {
        "description": "Get chat messages",
        "produces": [
          "application/json"
        ],
        "operationId": "GetChatMessages",
        "parameters": [
          {
            "type": "string",
            "name": "chatId",
            "in": "path",
            "required": true
          },
          {
            "type": "integer",
            "format": "int32",
            "default": 10,
            "name": "limit",
            "in": "query"
          },
          {
            "type": "integer",
            "format": "int64",
            "default": 0,
            "name": "before_timestamp",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "A list of chat messages",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/ChatMessage"
              }
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/error"
            }
          },
          "401": {
            "description": "unauthorized",
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
    },
    "/voice_messages": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Processes a voice message and returns a response.",
        "operationId": "SendVoiceMessage",
        "parameters": [
          {
            "description": "User message containing a voice message url",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "voiceMessageUrl": {
                  "description": "URL of the voice message file",
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
                "voiceMessageTranscript": {
                  "description": "Transcription of the input voice message.",
                  "type": "string"
                },
                "voiceMessageUrl": {
                  "description": "URL of the voice message input file.",
                  "type": "string"
                },
                "voiceResponseTranscript": {
                  "description": "Text transcription of the voice response.",
                  "type": "string"
                },
                "voiceResponseUrl": {
                  "description": "URL to the voice response file.",
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
          "401": {
            "description": "unauthorized",
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
    "ChatMessage": {
      "type": "object",
      "properties": {
        "curUsr": {
          "type": "boolean",
          "x-go-name": "IsFromCurrentUser"
        },
        "text": {
          "type": "string"
        },
        "timestamp": {
          "type": "integer",
          "format": "int64"
        },
        "userId": {
          "type": "string"
        }
      }
    },
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
  },
  "securityDefinitions": {
    "key": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header",
      "authorizationUrl": ""
    }
  },
  "security": [
    {
      "key": []
    }
  ]
}`))
}
