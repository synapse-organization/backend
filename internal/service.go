package internal

import (
	"barista/pkg/models"
	"github.com/gin-gonic/gin"
)

type Person struct {
	Name string
}

func (person Person) GetName(c *gin.Context) {
	c.JSON(200, gin.H{"name": person.Name})
}

func (person Person) PostName(c *gin.Context) {
	c.JSON(200, gin.H{"set name": person.Name})
}

func Run() {
	service := StartService()
	api := service.AddGroup("/api")

	service.AddRoutes(api, models.Route{
		Method: "GET",
		Path:   "/ping",
		Function: func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		},
	})

	a := Person{Name: "John Doe"}
	service.AddStructRoutes(api, a)

	service.Run(":8080")
}
