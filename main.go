package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

var db *gorm.DB
var err error

type User struct {
	gorm.Model
	Username string
	Password string
}

type Article struct {
	gorm.Model
	Title   string
	Content string
	UserID  uint
}

func signUpHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("err: %v", err),
		})
		return
	}
	db.Create(&User{Username: user.Username, Password: user.Password})
	c.JSON(http.StatusCreated, gin.H{
		"created": "ok",
	})
}

func signInHandler(c *gin.Context) {
	var user User
	var existUser User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("err: %v", err),
		})
		return
	}
	db.First(&existUser, user.ID)
	if user.Username == existUser.Username && user.Password == existUser.Password {
		c.JSON(http.StatusOK, gin.H{
			"signIn": "ok",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"signIn": "fail",
		})
	}
}

func createHandler(c *gin.Context) {
	var article Article
	if err := c.ShouldBind(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("err: %v", err),
		})
		return
	}
	db.Create(&Article{Title: article.Title, Content: article.Content, UserID: article.UserID})
	c.JSON(http.StatusCreated, gin.H{
		"created": "ok",
	})
}

func retrieveAllHandler(c *gin.Context) {
	var articles []Article
	db.Find(&articles, &Article{})
	c.JSON(http.StatusOK, gin.H{
		"result": articles,
	})
}

func retrieveHandler(c *gin.Context) {
	var article Article
	id := c.Param("id")
	db.First(&article, id)
	c.JSON(http.StatusOK, gin.H{
		"result": article,
	})
}

func updateHandler(c *gin.Context) {
	var article Article
	id := c.Param("id")
	fmt.Println(id)
	db.Where(id).Updates(article)
	c.JSON(http.StatusPartialContent, gin.H{
		"updated": "ok",
	})
}

func deleteHandler(c *gin.Context) {
	var article Article
	id := c.Param("id")
	db.Where(id).Delete(&article)
	c.JSON(http.StatusAccepted, gin.H{
		"deleted": "ok",
	})
}

func init() {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Seoul"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("[-] DB Connection Failed...")
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Article{})
}

func main() {
	boardServer().Run()
}

func boardServer() *gin.Engine {
	r := gin.Default()

	r.POST("/users/signup", signUpHandler)
	r.POST("/users/signin", signInHandler)
	r.GET("/articles", retrieveAllHandler)
	r.POST("/articles", createHandler)
	r.GET("/articles/:id", retrieveHandler)
	r.PUT("/articles/:id", updateHandler)
	r.DELETE("/articles/:id", deleteHandler)

	return r
}
