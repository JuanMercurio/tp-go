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
        "/cotizaciones": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Moneda"
                ],
                "summary": "Retorna las cotizaciones paginadas segun los filtros",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id de las monedas que queremos separados por espacios",
                        "name": "monedas",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "format": "date-time",
                        "description": "Fecha desde la cual quiero obtener cotizaciones (YYYY-MM-DD HH:MM:SS)",
                        "name": "fecha_inicial",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "format": "date-time",
                        "description": "Fecha hasta la cual quiero obtener cotizaciones (YYYY-MM-DD HH:MM:SS)",
                        "name": "fecha_final",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "El tamaño de las paginas, como maximo es 100, el default es 50",
                        "name": "tam_paginas",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Pagina a partir de la cual sera retornado el query",
                        "name": "pagina_inicial",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "La cantidad de paginas, como maximo es 10, el default es 1",
                        "name": "cant_paginas",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Ordenar segun alguno de estos valores: fecha(default), valor, nombre",
                        "name": "orden",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Indica si es ascendente o descendente, el default es desdencente",
                        "name": "orden_direccion",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Usuario elegido",
                        "name": "usuario",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Para incluir resumen indicar el valor debe ser si",
                        "name": "resumen",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "array",
                                "items": {}
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Moneda"
                ],
                "summary": "Llama para que se haga la cotizacion de las monedas",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Token de autorización",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "API que se utilizara para cotizar",
                        "name": "api",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "array",
                                "items": {}
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/monedas": {
            "get": {
                "description": "Obtiene una lista de todos las monedas disponibles.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Moneda"
                ],
                "summary": "Busca todas las monedas",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/ports.MonedaOutputDTO"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Si tenemos las credenciales podemos dar de alta una moneda",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Moneda"
                ],
                "summary": "Da de alta una moneda",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Token de autorización",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Simbolo de la moneda",
                        "name": "simbolo",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Nombre de la moneda nueva",
                        "name": "nombre",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        },
        "/usuarios": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Moneda"
                ],
                "summary": "Lista de usuario registrados",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/domain.Usuario"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Moneda"
                ],
                "summary": "Crear un usuario",
                "parameters": [
                    {
                        "type": "string",
                        "description": "nombre del nuevo usuario",
                        "name": "nombre",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Moneda"
                ],
                "summary": "Da de baja a un usuario",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "id del usuario que se desea dar de baja",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/usuarios/{id}/monedas": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Moneda"
                ],
                "summary": "Listar las monedas preferidas de un usuario",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "id del usuario",
                        "name": "id",
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
                                "$ref": "#/definitions/ports.MonedaOutputDTO"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.Usuario": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "nombre": {
                    "type": "string"
                }
            }
        },
        "ports.MonedaOutputDTO": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "nombre": {
                    "type": "string"
                },
                "simbolo": {
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
	Title:            "Criptomoneda API",
	Description:      "API en la cual podemos consultar cotizaciones de diferentes monedas",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
