package models

import "github.com/gin-gonic/gin"

type Route struct {
	Method   string
	Path     string
	Function gin.HandlerFunc
}
