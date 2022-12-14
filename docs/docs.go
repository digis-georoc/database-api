// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "DIGIS Project",
            "url": "https://www.uni-goettingen.de/de/643369.html",
            "email": "digis-info@uni-goettingen.de"
        },
        "license": {
            "name": "Data retrieved is licensed under CC BY-SA 4.0",
            "url": "https://creativecommons.org/licenses/by-sa/4.0/"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/geodata/sites": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "get site data in GeoJSON format",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sites"
                ],
                "summary": "Retrieve site data as GeoJSON",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "offset",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.GeoJSONFeatureCollection"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "Check connection to api",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "general"
                ],
                "summary": "Sample request",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/queries/authors": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "get authors",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "people"
                ],
                "summary": "Retrieve authors",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "offset",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.People"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/queries/authors/{personID}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "get authors by personID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "people"
                ],
                "summary": "Retrieve authors by personID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Person ID",
                        "name": "personID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.People"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/queries/citations": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "get citations",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "citations"
                ],
                "summary": "Retrieve citations",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "offset",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Citation"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/queries/citations/{citationID}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "get citations by citationID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "citations"
                ],
                "summary": "Retrieve citations by citationID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Citation ID",
                        "name": "citationID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Citation"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/queries/fulldata/{samplingfeatureid}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "get full dataset by samplingfeatureid",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "fulldata"
                ],
                "summary": "Retrieve full dataset by samplingfeatureid",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Samplingfeature identifier",
                        "name": "samplingfeatureid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.FullData"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/queries/samples": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get all samples matching the current filters\nMultiple values in a single filter must be comma separated",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "samples"
                ],
                "summary": "Retrieve all samples filtered by a variety of fields",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "offset",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "tectonic setting",
                        "name": "setting",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "location level 1",
                        "name": "location1",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "location level 2",
                        "name": "location2",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "location level 3",
                        "name": "location3",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "samplingfeature name",
                        "name": "samplename",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "sampling technique",
                        "name": "sampletech",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "land or sea",
                        "name": "landorsea",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "taxonomic classifier name",
                        "name": "rockclass",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "rock type",
                        "name": "rocktype",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "material",
                        "name": "material",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "chemical element",
                        "name": "majorelem",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Sample"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/queries/sites": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get all sites",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sites"
                ],
                "summary": "Retrieve all sites",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "offset",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Site"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/queries/sites/settings": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get all geological settings",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sites"
                ],
                "summary": "Retrieve all geological settings",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Site"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/queries/sites/{samplingfeatureID}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get sites by samplingfeatureID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sites"
                ],
                "summary": "Retrieve sites by samplingfeatureID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "samplingfeatureID",
                        "name": "samplingfeatureID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Site"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.Citation": {
            "type": "object",
            "properties": {
                "authors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.People"
                    }
                },
                "bookTitle": {
                    "type": "string"
                },
                "citationID": {
                    "type": "integer"
                },
                "citationLink": {
                    "type": "string"
                },
                "editors": {
                    "type": "string"
                },
                "firstPage": {
                    "type": "string"
                },
                "issue": {
                    "type": "string"
                },
                "journal": {
                    "type": "string"
                },
                "lastPage": {
                    "type": "string"
                },
                "publicationyear": {
                    "type": "integer"
                },
                "publisher": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "volume": {
                    "type": "string"
                }
            }
        },
        "model.FullData": {
            "type": "object",
            "properties": {
                "age_Max": {
                    "type": "integer"
                },
                "age_Min": {
                    "type": "integer"
                },
                "batches": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "comment": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "elevation_Max": {
                    "type": "integer"
                },
                "elevation_Min": {
                    "type": "integer"
                },
                "inclusion_Type": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "institution": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "string"
                        }
                    }
                },
                "item_Group": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "string"
                        }
                    }
                },
                "item_Name": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "string"
                        }
                    }
                },
                "land_Or_Sea": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "string"
                        }
                    }
                },
                "latitude": {
                    "type": "number"
                },
                "latitude_Max": {
                    "type": "string"
                },
                "latitude_Min": {
                    "type": "string"
                },
                "loc_Data": {},
                "location_Names": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "string"
                        }
                    }
                },
                "location_Num": {
                    "type": "integer"
                },
                "location_Types": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "string"
                        }
                    }
                },
                "longitude": {
                    "type": "number"
                },
                "longitude_Max": {
                    "type": "string"
                },
                "longitude_Min": {
                    "type": "string"
                },
                "material": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "method": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "mineral": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "references": {},
                "rock_Class": {
                    "type": "string"
                },
                "rock_Texture": {
                    "type": "string"
                },
                "rock_Type": {
                    "type": "string"
                },
                "sampleIds": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "sample_Num": {
                    "type": "integer"
                },
                "standard_Names": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "string"
                        }
                    }
                },
                "standard_Values": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "number"
                        }
                    }
                },
                "tectonic_Setting": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "string"
                        }
                    }
                },
                "unique_id": {
                    "type": "string"
                },
                "units": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "string"
                        }
                    }
                },
                "values": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "number"
                        }
                    }
                }
            }
        },
        "model.GeoJSONFeature": {
            "type": "object",
            "properties": {
                "geometry": {
                    "$ref": "#/definitions/model.Geometry"
                },
                "id": {
                    "type": "string"
                },
                "properties": {
                    "type": "object",
                    "additionalProperties": true
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "model.GeoJSONFeatureCollection": {
            "type": "object",
            "properties": {
                "features": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.GeoJSONFeature"
                    }
                },
                "numberMatched": {
                    "type": "integer"
                },
                "numberReturned": {
                    "type": "integer"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "model.Geometry": {
            "type": "object",
            "properties": {
                "coordinates": {
                    "type": "array",
                    "items": {}
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "model.People": {
            "type": "object",
            "properties": {
                "personFirstName": {
                    "type": "string"
                },
                "personId": {
                    "type": "integer"
                },
                "personLastName": {
                    "type": "string"
                }
            }
        },
        "model.Sample": {
            "type": "object",
            "properties": {
                "elevationPrecision": {
                    "type": "number"
                },
                "elevationPrecisionComment": {
                    "type": "string"
                },
                "samplingFeatureCode": {
                    "type": "string"
                },
                "samplingFeatureDescription": {
                    "type": "string"
                },
                "samplingFeatureID": {
                    "type": "integer"
                },
                "samplingFeatureName": {
                    "type": "string"
                },
                "samplingFeatureTypeCV": {
                    "type": "string"
                },
                "samplingFeatureUUID": {
                    "type": "integer"
                }
            }
        },
        "model.Site": {
            "type": "object",
            "properties": {
                "latitude": {
                    "type": "number"
                },
                "locationPrecision": {
                    "type": "number"
                },
                "locationPrecisionComment": {
                    "type": "string"
                },
                "longitude": {
                    "type": "number"
                },
                "samplingFeatureID": {
                    "type": "integer"
                },
                "setting": {
                    "type": "string"
                },
                "siteDescription": {
                    "type": "string"
                },
                "siteTypeCV": {
                    "type": "string"
                },
                "spatialReferenceID": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "description": "Accesskey based security scheme to secure api groups \"/queries/*\" and \"/geodata/*\"",
            "type": "apiKey",
            "name": "DIGIS-API-ACCESSKEY",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1.0",
	Host:             "api-test.georoc.eu",
	BasePath:         "/api/v1",
	Schemes:          []string{"https", "http"},
	Title:            "DIGIS Database API",
	Description:      "This is the database api for the new GeoROC datamodel\n\nNote: Semicolon (;) in queries are not allowed and need to be url-encoded as per this issue: golang.org/issue/25192",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
