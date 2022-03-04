package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"text/template"
)

var (
	db  *gorm.DB
	err error
	tpl *template.Template
)

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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func signUpHandler(c *gin.Context) {
	var user User
	var existUser User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("err: %v", err),
		})
		return
	}
	if db.First(&existUser, "username", user.Username) == nil {
		hash, _ := HashPassword(user.Password)
		db.Create(&User{Username: user.Username, Password: hash})
		c.JSON(http.StatusCreated, gin.H{
			"created": "ok",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("err: %v", err),
		})
	}
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
	db.First(&existUser, "username", user.Username)
	if CheckPasswordHash(user.Password, existUser.Password) {
		c.JSON(http.StatusOK, gin.H{
			"accessToken":  "",
			"refreshToken": "",
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

func indexHandler(c *gin.Context) {
	tpl.ExecuteTemplate(c.Writer, "index.gohtml", nil)
}

func init() {
	tpl = template.Must(template.ParseGlob("pages/*.gohtml"))
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

	r.GET("/api/", indexHandler)
	r.POST("/api/users/signup", signUpHandler)
	r.POST("/api/users/signin", signInHandler)
	r.GET("/api/articles", retrieveAllHandler)
	r.POST("/api/articles", createHandler)
	r.GET("/api/articles/:id", retrieveHandler)
	r.PUT("/api/articles/:id", updateHandler)
	r.DELETE("/api/articles/:id", deleteHandler)

	return r
}
