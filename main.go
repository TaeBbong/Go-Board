package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"utils"
)

func init() {
	utils.SetupDB()
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
