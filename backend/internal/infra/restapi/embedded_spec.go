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
        "summary": "Sends a text message to the AI and receives a response.",
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
                "chat": {
                  "description": "If the chatId is not provided in the request, the new chat will be created and this chat will be returned in the response.",
                  "type": "object",
                  "$ref": "#/definitions/Chat"
                },
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
            "description": "timestamp starting from which messages are to be fetched",
            "name": "timestamp",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "A list of chat messages",
            "schema": {
              "type": "object",
              "properties": {
                "messages": {
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/ChatMessage"
                  }
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
    "/chat_voice_messages": {
      "post": {
        "description": "This endpoint receives a user's audio message and returns the AI's response.",
        "consumes": [
          "multipart/form-data"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Sends an audio message to the AI and receives a response.",
        "operationId": "SendVoiceMessage",
        "parameters": [
          {
            "type": "file",
            "description": "The audio file to be sent.",
            "name": "file",
            "in": "formData",
            "required": true
          },
          {
            "type": "string",
            "description": "The chat ID.",
            "name": "chatId",
            "in": "formData"
          },
          {
            "type": "integer",
            "format": "int64",
            "description": "The timestamp of the message.",
            "name": "timestamp",
            "in": "formData",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "type": "object",
              "properties": {
                "chat": {
                  "description": "If the chatId is not provided in the request, the new chat will be created and this chat will be returned in the response.",
                  "type": "object",
                  "$ref": "#/definitions/Chat"
                },
                "replyAudio": {
                  "description": "Url to the AI's audio response.",
                  "type": "string"
                },
                "replyTime": {
                  "description": "The timestamp of the AI's response.",
                  "type": "integer",
                  "format": "int64"
                },
                "replyTxt": {
                  "description": "AI's text response to the user's message.",
                  "type": "string"
                },
                "usrAudio": {
                  "description": "Url to the user's audio message.",
                  "type": "string"
                },
                "usrTime": {
                  "description": "The timestamp of the user's message.",
                  "type": "integer",
                  "format": "int64"
                },
                "usrTxt": {
                  "description": "The user's message in text format.",
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
    },
    "/chats": {
      "get": {
        "description": "Get chats",
        "produces": [
          "application/json"
        ],
        "operationId": "GetChats",
        "parameters": [
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
            "description": "timestamp starting from which chats are to be fetched",
            "name": "timestamp",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "A list of chats",
            "schema": {
              "type": "object",
              "properties": {
                "chats": {
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/Chat"
                  }
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
      },
      "post": {
        "description": "This endpoint creates a new chat.",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Creates a new chat.",
        "operationId": "CreateChat",
        "parameters": [
          {
            "description": "Chat information",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "typ",
                "time"
              ],
              "properties": {
                "time": {
                  "description": "The timestamp of the chat.",
                  "type": "integer",
                  "format": "int64"
                },
                "typ": {
                  "x-go-name": "ChatType",
                  "$ref": "#/definitions/ChatType"
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
                "chat": {
                  "description": "The created chat.",
                  "type": "object",
                  "$ref": "#/definitions/Chat"
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
    "/go-to-message": {
      "post": {
        "description": "This endpoint goes to a specific prepared message in the chat.",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Goes to a specific prepared message in the chat.",
        "operationId": "GoToMessage",
        "parameters": [
          {
            "description": "Information about the chat and the message index",
            "name": "body",
            "in": "body",
            "schema": {
              "type": "object",
              "required": [
                "chatId",
                "msgIdx"
              ],
              "properties": {
                "chatId": {
                  "description": "The chat ID.",
                  "type": "string"
                },
                "msgIdx": {
                  "description": "The index of the message.",
                  "type": "integer",
                  "format": "int32"
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
                "msg": {
                  "$ref": "#/definitions/ChatMessage"
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
    "Chat": {
      "type": "object",
      "properties": {
        "chatId": {
          "type": "string"
        },
        "cur_q": {
          "description": "The current question in the chat. 0 based index.",
          "type": "integer",
          "format": "int32",
          "x-go-name": "CurrentQuestionIDx"
        },
        "prep_msgs": {
          "type": "array",
          "items": {
            "type": "string"
          },
          "x-go-name": "PreparedMessages"
        },
        "time": {
          "type": "integer",
          "format": "int64"
        },
        "title": {
          "type": "string"
        },
        "typ": {
          "$ref": "#/definitions/ChatType"
        }
      }
    },
    "ChatMessage": {
      "type": "object",
      "properties": {
        "audio": {
          "type": "string",
          "x-go-name": "AudioURL"
        },
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
    "ChatType": {
      "type": "integer",
      "format": "int32",
      "enum": [
        1,
        2
      ],
      "x-enum-descriptions": [
        "General",
        "Training specific questions from job interviews"
      ]
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
        "summary": "Sends a text message to the AI and receives a response.",
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
                "chat": {
                  "description": "If the chatId is not provided in the request, the new chat will be created and this chat will be returned in the response.",
                  "type": "object",
                  "$ref": "#/definitions/Chat"
                },
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
            "description": "timestamp starting from which messages are to be fetched",
            "name": "timestamp",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "A list of chat messages",
            "schema": {
              "type": "object",
              "properties": {
                "messages": {
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/ChatMessage"
                  }
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
    "/chat_voice_messages": {
      "post": {
        "description": "This endpoint receives a user's audio message and returns the AI's response.",
        "consumes": [
          "multipart/form-data"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Sends an audio message to the AI and receives a response.",
        "operationId": "SendVoiceMessage",
        "parameters": [
          {
            "type": "file",
            "description": "The audio file to be sent.",
            "name": "file",
            "in": "formData",
            "required": true
          },
          {
            "type": "string",
            "description": "The chat ID.",
            "name": "chatId",
            "in": "formData"
          },
          {
            "type": "integer",
            "format": "int64",
            "description": "The timestamp of the message.",
            "name": "timestamp",
            "in": "formData",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "type": "object",
              "properties": {
                "chat": {
                  "description": "If the chatId is not provided in the request, the new chat will be created and this chat will be returned in the response.",
                  "type": "object",
                  "$ref": "#/definitions/Chat"
                },
                "replyAudio": {
                  "description": "Url to the AI's audio response.",
                  "type": "string"
                },
                "replyTime": {
                  "description": "The timestamp of the AI's response.",
                  "type": "integer",
                  "format": "int64"
                },
                "replyTxt": {
                  "description": "AI's text response to the user's message.",
                  "type": "string"
                },
                "usrAudio": {
                  "description": "Url to the user's audio message.",
                  "type": "string"
                },
                "usrTime": {
                  "description": "The timestamp of the user's message.",
                  "type": "integer",
                  "format": "int64"
                },
                "usrTxt": {
                  "description": "The user's message in text format.",
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
    },
    "/chats": {
      "get": {
        "description": "Get chats",
        "produces": [
          "application/json"
        ],
        "operationId": "GetChats",
        "parameters": [
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
            "description": "timestamp starting from which chats are to be fetched",
            "name": "timestamp",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "A list of chats",
            "schema": {
              "type": "object",
              "properties": {
                "chats": {
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/Chat"
                  }
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
      },
      "post": {
        "description": "This endpoint creates a new chat.",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Creates a new chat.",
        "operationId": "CreateChat",
        "parameters": [
          {
            "description": "Chat information",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "typ",
                "time"
              ],
              "properties": {
                "time": {
                  "description": "The timestamp of the chat.",
                  "type": "integer",
                  "format": "int64"
                },
                "typ": {
                  "x-go-name": "ChatType",
                  "$ref": "#/definitions/ChatType"
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
                "chat": {
                  "description": "The created chat.",
                  "type": "object",
                  "$ref": "#/definitions/Chat"
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
    "/go-to-message": {
      "post": {
        "description": "This endpoint goes to a specific prepared message in the chat.",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Goes to a specific prepared message in the chat.",
        "operationId": "GoToMessage",
        "parameters": [
          {
            "description": "Information about the chat and the message index",
            "name": "body",
            "in": "body",
            "schema": {
              "type": "object",
              "required": [
                "chatId",
                "msgIdx"
              ],
              "properties": {
                "chatId": {
                  "description": "The chat ID.",
                  "type": "string"
                },
                "msgIdx": {
                  "description": "The index of the message.",
                  "type": "integer",
                  "format": "int32"
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
                "msg": {
                  "$ref": "#/definitions/ChatMessage"
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
    "Chat": {
      "type": "object",
      "properties": {
        "chatId": {
          "type": "string"
        },
        "cur_q": {
          "description": "The current question in the chat. 0 based index.",
          "type": "integer",
          "format": "int32",
          "x-go-name": "CurrentQuestionIDx"
        },
        "prep_msgs": {
          "type": "array",
          "items": {
            "type": "string"
          },
          "x-go-name": "PreparedMessages"
        },
        "time": {
          "type": "integer",
          "format": "int64"
        },
        "title": {
          "type": "string"
        },
        "typ": {
          "$ref": "#/definitions/ChatType"
        }
      }
    },
    "ChatMessage": {
      "type": "object",
      "properties": {
        "audio": {
          "type": "string",
          "x-go-name": "AudioURL"
        },
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
    "ChatType": {
      "type": "integer",
      "format": "int32",
      "enum": [
        1,
        2
      ],
      "x-enum-descriptions": [
        "General",
        "Training specific questions from job interviews"
      ]
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
