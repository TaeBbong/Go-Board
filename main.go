package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"net/http"
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

func indexPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"content": "This is IndexPage",
		"isLogin": false,
	})
}

func postListPageHandler(c *gin.Context) {
	var articles []Article
	db.Find(&articles, &Article{})
	c.HTML(http.StatusOK, "post_list.html", gin.H{
		"posts": articles,
	})
}

func postDetailPageHandler(c *gin.Context) {
	var article Article
	id := c.Param("id")
	db.First(&article, id)
	c.HTML(http.StatusOK, "post_detail.html", gin.H{
		"post": article,
	})
}

func postCreatePageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "post_create.html", gin.H{})
}

func postUpdatePageHandler(c *gin.Context) {
	var article Article
	id := c.Param("id")
	db.First(&article, id)
	c.HTML(http.StatusOK, "post_edit.html", gin.H{
		"post": article,
	})
}

func userSignUpPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "user_signup.html", gin.H{})
}

func userSignInPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "user_signin.html", gin.H{})
}

func postCreateHandler(c *gin.Context) {
	var article Article
	if err := c.ShouldBind(&article); err != nil { // TODO: form 관련 사용법 체크
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("err: %v", err),
		})
		return
	}
	db.Create(&Article{Title: article.Title, Content: article.Content, UserID: article.UserID})
	c.Redirect(http.StatusFound, "/post")
}

func postUpdateHandler(c *gin.Context) {
	var article Article
	if err := c.ShouldBind(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("err: %v", err),
		})
		return
	}
	id := c.Param("id")
	db.Where(id).Updates(&Article{Title: article.Title, Content: article.Content})
	c.Redirect(http.StatusFound, "/post")
}

func postDeleteHandler(c *gin.Context) {
	var article Article
	id := c.Param("id")
	db.Where(id).Delete(&article)
	c.Redirect(http.StatusAccepted, "/post")
}

func userSignUpHandler(c *gin.Context) {
	var user User
	var existUser User
	if err := c.ShouldBind(&user); err != nil {
		c.HTML(http.StatusBadRequest, "user_signup_error.html", gin.H{})
		return
	}
	if db.First(&existUser, "username", user.Username) == nil {
		hash, _ := HashPassword(user.Password)
		db.Create(&User{Username: user.Username, Password: hash})
		c.Redirect(http.StatusCreated, "/user/signin")
	} else {
		c.HTML(http.StatusBadRequest, "user_signup_error.html", gin.H{})
	}
}

func userSignInHandler(c *gin.Context) {
	var user User
	var existUser User
	if err := c.ShouldBind(&user); err != nil {
		c.HTML(http.StatusBadRequest, "user_signin_error.html", gin.H{})
		return
	}
	db.First(&existUser, "username", user.Username)
	if CheckPasswordHash(user.Password, existUser.Password) {
		c.Redirect(http.StatusOK, "/")
	} else {
		c.HTML(http.StatusBadRequest, "user_signin_error.html", gin.H{})
	}
}

func userSignOutHandler(c *gin.Context) {
	c.Redirect(http.StatusOK, "/")
}

func init() {
	tpl = template.Must(template.ParseGlob("pages/*.html"))
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
	// boardAPIServer().Run()
}

func boardServer() *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("pages/*")

	r.GET("/", indexPageHandler)
	r.GET("/post", postListPageHandler)
	r.GET("/post/:id", postDetailPageHandler)
	r.GET("/post/create", postCreatePageHandler)
	r.POST("/post/create", postCreateHandler)
	r.GET("/post/:id/edit", postUpdatePageHandler)
	r.POST("/post/:id/edit", postUpdateHandler)
	r.GET("/post/:id/delete", postDeleteHandler)

	r.GET("/user/signup", userSignUpPageHandler)
	r.POST("/user/signup", userSignUpHandler)
	r.GET("/user/signin", userSignInPageHandler)
	r.POST("/user/signin", userSignInHandler)
	r.GET("/user/signout", userSignOutHandler)

	return r
}

func boardAPIServer() *gin.Engine {
	r := gin.Default()

	// api version
	r.POST("/api/users/signup", signUpHandler)
	r.POST("/api/users/signin", signInHandler)
	r.GET("/api/articles", retrieveAllHandler)
	r.POST("/api/articles", createHandler)
	r.GET("/api/articles/:id", retrieveHandler)
	r.PUT("/api/articles/:id", updateHandler)
	r.DELETE("/api/articles/:id", deleteHandler)

	return r
}
