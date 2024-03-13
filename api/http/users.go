package http

import "github.com/gin-gonic/gin"

type User struct {
}

func (userHandler User) GetUsers(c *gin.Context) {
	c.JSON(200, gin.H{"users": "get"})
}

func (userHandler User) GetSingIn(c *gin.Context) {
	c.JSON(200, gin.H{"users": "singin"})
}

func (userHandler User) PostSignUp(c *gin.Context) {
	c.JSON(200, gin.H{"users": "signup"})
}
