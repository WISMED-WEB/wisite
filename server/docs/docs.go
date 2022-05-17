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
            "name": "API Support"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/admin/activate": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "activate or deactivate a user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "unique user name",
                        "name": "uname",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "true: activate, false: deactivate",
                        "name": "flag",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - action successfully"
                    },
                    "400": {
                        "description": "Fail - invalid true/false flag"
                    },
                    "401": {
                        "description": "Fail - unauthorized error"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/admin/officialize": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "officialize or un-officialize a user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "unique user name",
                        "name": "uname",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "true: officialize, false: un-officialize",
                        "name": "flag",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - action successfully"
                    },
                    "400": {
                        "description": "Fail - invalid true/false flag"
                    },
                    "401": {
                        "description": "Fail - unauthorized error"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/admin/onlines": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "get all online users",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user filter with uname wildcard(*)",
                        "name": "uname",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - list successfully"
                    },
                    "401": {
                        "description": "Fail - unauthorized error"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/admin/spa/menu": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "get tailored side menu for different user group",
                "responses": {
                    "200": {
                        "description": "OK - get menu successfully"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/admin/users": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "get all users' info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user filter with uname wildcard(*)",
                        "name": "uname",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "user filter with name wildcard(*)",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "user filter with active status",
                        "name": "active",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - list successfully"
                    },
                    "401": {
                        "description": "Fail - unauthorized error"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/file/fileitem": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "file"
                ],
                "summary": "get fileitems by given path or id.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "path to a file",
                        "name": "path",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "file's id",
                        "name": "id",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - get fileitems successfully"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/file/pathcontent": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "file"
                ],
                "summary": "get content under specific path.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "path to some level",
                        "name": "path",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - upload successfully"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/file/upload": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "file"
                ],
                "summary": "upload file action.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "note for uploading file",
                        "name": "note",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "1st category for uploading file",
                        "name": "group0",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "2nd category for uploading file",
                        "name": "group1",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "3rd category for uploading file",
                        "name": "group2",
                        "in": "formData"
                    },
                    {
                        "type": "file",
                        "description": "file path for uploading",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - upload successfully"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/post/ids": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Post"
                ],
                "summary": "get a batch of Post id group.",
                "responses": {
                    "200": {
                        "description": "OK - get successfully"
                    },
                    "400": {
                        "description": "Fail - "
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/post/template": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Post"
                ],
                "summary": "get Post template for dev reference.",
                "responses": {
                    "200": {
                        "description": "OK - upload successfully"
                    }
                }
            }
        },
        "/api/post/upload": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Post"
                ],
                "summary": "upload a Post by filling a Post template.",
                "responses": {
                    "200": {
                        "description": "OK - upload successfully"
                    },
                    "400": {
                        "description": "Fail - incorrect Post format"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/rel/action/{whom}": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "relation"
                ],
                "summary": "relation actions",
                "parameters": [
                    {
                        "type": "string",
                        "description": "which action to apply, accept [follow, unfollow, block, unblock, mute, unmute]",
                        "name": "action",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "whose uname you want to follow",
                        "name": "whom",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - following successfully"
                    },
                    "400": {
                        "description": "Fail - invalid action type"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/rel/content/{type}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "relation"
                ],
                "summary": "get all relation users for one type",
                "parameters": [
                    {
                        "type": "string",
                        "description": "relation content type to apply, accept [following, follower, blocked, muted]",
                        "name": "type",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - got following successfully"
                    },
                    "400": {
                        "description": "Fail - invalid relation content type"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/sign-out/": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sign"
                ],
                "summary": "sign out action.",
                "responses": {
                    "200": {
                        "description": "OK - sign-out successfully"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/sign/in": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sign"
                ],
                "summary": "sign in action. if ok, got token",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user name or email",
                        "name": "uname",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "format": "password",
                        "description": "password",
                        "name": "pwd",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - sign-in successfully"
                    },
                    "400": {
                        "description": "Fail - incorrect password"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/sign/new": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sign"
                ],
                "summary": "sign up action, step 1. send user's basic info for registry",
                "parameters": [
                    {
                        "type": "string",
                        "description": "unique user name",
                        "name": "uname",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "format": "email",
                        "description": "user's email",
                        "name": "email",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "user's real full name",
                        "name": "name",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "user's password",
                        "name": "pwd",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - then waiting for verification code"
                    },
                    "400": {
                        "description": "Fail - invalid registry fields"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/sign/reset-pwd": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sign"
                ],
                "summary": "reset password action, step 1. send verification code to user's email for authentication",
                "parameters": [
                    {
                        "type": "string",
                        "description": "unique user name",
                        "name": "uname",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "format": "email",
                        "description": "user's email",
                        "name": "email",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - then waiting for verification code"
                    },
                    "400": {
                        "description": "Fail - invalid registry fields"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/sign/verify-email": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sign"
                ],
                "summary": "sign up action, step 2. send back email verification code",
                "parameters": [
                    {
                        "type": "string",
                        "description": "unique user name",
                        "name": "uname",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "verification code (in user's email)",
                        "name": "code",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - sign-up successfully"
                    },
                    "400": {
                        "description": "Fail - incorrect verification code"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/sign/verify-reset-pwd": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "sign"
                ],
                "summary": "reset password action, step 2. send back verification code for updating password",
                "parameters": [
                    {
                        "type": "string",
                        "description": "unique user name",
                        "name": "uname",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "verification code (in user's email)",
                        "name": "code",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "new password",
                        "name": "pwd",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK   - password updated successfully"
                    },
                    "400": {
                        "description": "Fail - incorrect verification code"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/system/ver": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "system"
                ],
                "summary": "get this api service version",
                "responses": {
                    "200": {
                        "description": "OK - get its version"
                    }
                }
            }
        },
        "/api/system/ver-tag": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "system"
                ],
                "summary": "get this api service project github version tag",
                "responses": {
                    "200": {
                        "description": "OK - get its tag"
                    }
                }
            }
        },
        "/api/user/avatar": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "get user avatar src as base64",
                "responses": {
                    "200": {
                        "description": "OK - get avatar src base64"
                    },
                    "404": {
                        "description": "Fail - avatar is empty"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/user/heartbeats": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "frequently call this to indicate that front-end user is active.",
                "responses": {
                    "200": {
                        "description": "OK - heartbeats successfully"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/user/profile": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "get user profile",
                "responses": {
                    "200": {
                        "description": "OK - profile get successfully"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/user/setprofile": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "set user profile",
                "parameters": [
                    {
                        "type": "string",
                        "description": "phone number",
                        "name": "phone",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "address",
                        "name": "addr",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "city",
                        "name": "city",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "country",
                        "name": "country",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "personal id type",
                        "name": "pidtype",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "personal id",
                        "name": "pid",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "gender M/F",
                        "name": "gender",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "date of birth",
                        "name": "dob",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "job position",
                        "name": "position",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "title",
                        "name": "title",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "employer",
                        "name": "employer",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "biography",
                        "name": "bio",
                        "in": "formData"
                    },
                    {
                        "type": "file",
                        "description": "avatar",
                        "name": "avatar",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - profile set successfully"
                    },
                    "400": {
                        "description": "Fail - invalid set fields"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "127.0.0.1:1323",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "WISMED WISITE API",
	Description:      "This is wismed wisite-api server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
