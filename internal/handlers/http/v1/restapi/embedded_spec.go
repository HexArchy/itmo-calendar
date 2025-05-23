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
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Itmo calendar sync service",
    "title": "Itmo Calendar",
    "version": "1.0.0"
  },
  "basePath": "/api/v1",
  "paths": {
    "/health": {
      "get": {
        "security": [],
        "description": "Verifies the API is operational and returns its status.",
        "tags": [
          "System"
        ],
        "summary": "Health check endpoint.",
        "operationId": "healthCheck",
        "responses": {
          "200": {
            "description": "Service is healthy."
          },
          "503": {
            "description": "Service is unhealthy.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/subscribe": {
      "post": {
        "description": "Subscribes user by ISU and password, generates and stores iCal file.",
        "tags": [
          "CalDav"
        ],
        "summary": "Subscribe and generate iCal for user.",
        "operationId": "subscribeSchedule",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/SubscribeRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Subscription successful.",
            "schema": {
              "$ref": "#/definitions/SubscribeResponse"
            }
          },
          "400": {
            "description": "Bad request.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "description": "Internal server error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/{isu}/ical": {
      "get": {
        "description": "Returns the iCalendar (.ics) file for the user with the given ISU.",
        "produces": [
          "text/calendar"
        ],
        "tags": [
          "CalDav"
        ],
        "summary": "Get user's iCal file by ISU.",
        "operationId": "getICal",
        "parameters": [
          {
            "type": "integer",
            "format": "int64",
            "description": "ISU of the user.",
            "name": "isu",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "iCal file.",
            "schema": {
              "type": "string",
              "format": "binary"
            }
          },
          "404": {
            "description": "Not found.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "description": "Internal server error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/{isu}/schedule": {
      "get": {
        "description": "Returns the schedule for the user with the given ISU.",
        "tags": [
          "Schedule"
        ],
        "summary": "Get user's schedule by ISU.",
        "operationId": "getSchedule",
        "parameters": [
          {
            "type": "integer",
            "format": "int64",
            "description": "ISU of the user.",
            "name": "isu",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "User's schedule.",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/ScheduleItem"
              }
            }
          },
          "404": {
            "description": "Not found.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "description": "Internal server error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Error": {
      "type": "object",
      "properties": {
        "error": {
          "type": "string",
          "example": "Service Unavailable"
        },
        "message": {
          "type": "string",
          "example": "The service is currently unavailable. Please try again later."
        }
      }
    },
    "ScheduleItem": {
      "type": "object",
      "required": [
        "date",
        "lessons"
      ],
      "properties": {
        "date": {
          "type": "string",
          "format": "date",
          "example": "2024-06-01"
        },
        "lessons": {
          "type": "array",
          "items": {
            "type": "object",
            "required": [
              "subject",
              "type",
              "teacher_name",
              "room",
              "building",
              "format",
              "group",
              "time_start",
              "time_end"
            ],
            "properties": {
              "building": {
                "type": "string",
                "example": "Main"
              },
              "format": {
                "type": "string",
                "example": "Offline"
              },
              "group": {
                "type": "string",
                "example": "A1"
              },
              "note": {
                "type": "string",
                "example": "Bring calculator"
              },
              "room": {
                "type": "string",
                "example": "101"
              },
              "subject": {
                "type": "string",
                "example": "Mathematics"
              },
              "teacher_name": {
                "type": "string",
                "example": "Dr. Ivanov"
              },
              "time_end": {
                "type": "string",
                "format": "date-time",
                "example": "2024-06-01T10:30:00Z"
              },
              "time_start": {
                "type": "string",
                "format": "date-time",
                "example": "2024-06-01T09:00:00Z"
              },
              "type": {
                "type": "string",
                "example": "Lecture"
              },
              "zoom_url": {
                "type": "string",
                "example": "https://zoom.us/j/123456789"
              }
            }
          }
        }
      }
    },
    "SubscribeRequest": {
      "type": "object",
      "required": [
        "isu",
        "password"
      ],
      "properties": {
        "isu": {
          "type": "integer",
          "format": "int64",
          "example": 123456789
        },
        "password": {
          "type": "string",
          "example": "user_password"
        }
      }
    },
    "SubscribeResponse": {
      "type": "object",
      "properties": {
        "message": {
          "type": "string",
          "example": "Subscription successful. iCal generated."
        }
      }
    }
  },
  "securityDefinitions": {
    "JWT": {
      "description": "JWT token for user authentication",
      "type": "apiKey",
      "name": "X-Auth-Token",
      "in": "header"
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Itmo calendar sync service",
    "title": "Itmo Calendar",
    "version": "1.0.0"
  },
  "basePath": "/api/v1",
  "paths": {
    "/health": {
      "get": {
        "security": [],
        "description": "Verifies the API is operational and returns its status.",
        "tags": [
          "System"
        ],
        "summary": "Health check endpoint.",
        "operationId": "healthCheck",
        "responses": {
          "200": {
            "description": "Service is healthy."
          },
          "503": {
            "description": "Service is unhealthy.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/subscribe": {
      "post": {
        "description": "Subscribes user by ISU and password, generates and stores iCal file.",
        "tags": [
          "CalDav"
        ],
        "summary": "Subscribe and generate iCal for user.",
        "operationId": "subscribeSchedule",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/SubscribeRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Subscription successful.",
            "schema": {
              "$ref": "#/definitions/SubscribeResponse"
            }
          },
          "400": {
            "description": "Bad request.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "description": "Internal server error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/{isu}/ical": {
      "get": {
        "description": "Returns the iCalendar (.ics) file for the user with the given ISU.",
        "produces": [
          "text/calendar"
        ],
        "tags": [
          "CalDav"
        ],
        "summary": "Get user's iCal file by ISU.",
        "operationId": "getICal",
        "parameters": [
          {
            "type": "integer",
            "format": "int64",
            "description": "ISU of the user.",
            "name": "isu",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "iCal file.",
            "schema": {
              "type": "string",
              "format": "binary"
            }
          },
          "404": {
            "description": "Not found.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "description": "Internal server error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/{isu}/schedule": {
      "get": {
        "description": "Returns the schedule for the user with the given ISU.",
        "tags": [
          "Schedule"
        ],
        "summary": "Get user's schedule by ISU.",
        "operationId": "getSchedule",
        "parameters": [
          {
            "type": "integer",
            "format": "int64",
            "description": "ISU of the user.",
            "name": "isu",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "User's schedule.",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/ScheduleItem"
              }
            }
          },
          "404": {
            "description": "Not found.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "500": {
            "description": "Internal server error.",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Error": {
      "type": "object",
      "properties": {
        "error": {
          "type": "string",
          "example": "Service Unavailable"
        },
        "message": {
          "type": "string",
          "example": "The service is currently unavailable. Please try again later."
        }
      }
    },
    "ScheduleItem": {
      "type": "object",
      "required": [
        "date",
        "lessons"
      ],
      "properties": {
        "date": {
          "type": "string",
          "format": "date",
          "example": "2024-06-01"
        },
        "lessons": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/ScheduleItemLessonsItems0"
          }
        }
      }
    },
    "ScheduleItemLessonsItems0": {
      "type": "object",
      "required": [
        "subject",
        "type",
        "teacher_name",
        "room",
        "building",
        "format",
        "group",
        "time_start",
        "time_end"
      ],
      "properties": {
        "building": {
          "type": "string",
          "example": "Main"
        },
        "format": {
          "type": "string",
          "example": "Offline"
        },
        "group": {
          "type": "string",
          "example": "A1"
        },
        "note": {
          "type": "string",
          "example": "Bring calculator"
        },
        "room": {
          "type": "string",
          "example": "101"
        },
        "subject": {
          "type": "string",
          "example": "Mathematics"
        },
        "teacher_name": {
          "type": "string",
          "example": "Dr. Ivanov"
        },
        "time_end": {
          "type": "string",
          "format": "date-time",
          "example": "2024-06-01T10:30:00Z"
        },
        "time_start": {
          "type": "string",
          "format": "date-time",
          "example": "2024-06-01T09:00:00Z"
        },
        "type": {
          "type": "string",
          "example": "Lecture"
        },
        "zoom_url": {
          "type": "string",
          "example": "https://zoom.us/j/123456789"
        }
      }
    },
    "SubscribeRequest": {
      "type": "object",
      "required": [
        "isu",
        "password"
      ],
      "properties": {
        "isu": {
          "type": "integer",
          "format": "int64",
          "example": 123456789
        },
        "password": {
          "type": "string",
          "example": "user_password"
        }
      }
    },
    "SubscribeResponse": {
      "type": "object",
      "properties": {
        "message": {
          "type": "string",
          "example": "Subscription successful. iCal generated."
        }
      }
    }
  },
  "securityDefinitions": {
    "JWT": {
      "description": "JWT token for user authentication",
      "type": "apiKey",
      "name": "X-Auth-Token",
      "in": "header"
    }
  }
}`))
}
