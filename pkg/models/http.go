package models

import "github.com/gin-gonic/gin"

type Route struct {
	Method   string
	Path     string
	Function gin.HandlerFunc
}

type RequestMethod string

const (
	POST RequestMethod = "POST"
	GET  RequestMethod = "GET"
	PUT RequestMethod = "PUT"
	PATCH RequestMethod = "PATCH"
	DELETE RequestMethod = "DELETE"
)
