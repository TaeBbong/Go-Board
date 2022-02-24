package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"models"
)

func init() {
	models.SetupDB()
	fmt.Println("init from main.go")
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
