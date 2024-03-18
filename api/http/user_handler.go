package http

import "github.com/gin-gonic/gin"

type User struct {
}

func (u User) GetUser(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
