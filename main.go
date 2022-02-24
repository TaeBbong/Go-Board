package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"models"
	"services"
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
	handlers := func(c *gin.Context) {
		services.GetAll()
	}
	r.GET("/articles", handlers)
	r.Run()
}
