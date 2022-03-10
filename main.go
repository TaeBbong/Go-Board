package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	db  *gorm.DB
	err error
	tpl *template.Template
	fs  http.FileSystem = http.Dir("./upload")
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

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.Next()
}

func indexPageHandler(c *gin.Context) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"content":  "This is IndexPage",
		"username": session.Get("username"),
		"userID":   session.Get("userID"),
	})
}

func postListPageHandler(c *gin.Context) {
	session := sessions.Default(c)
	var articles []Article
	db.Find(&articles, &Article{})
	c.HTML(http.StatusOK, "post_list.html", gin.H{
		"posts":    articles,
		"username": session.Get("username"),
		"userID":   session.Get("userID"),
	})
}

func postDetailPageHandler(c *gin.Context) {
	session := sessions.Default(c)
	var article Article
	id := c.Param("id")
	db.First(&article, id)
	c.HTML(http.StatusOK, "post_detail.html", gin.H{
		"post":     article,
		"username": session.Get("username"),
		"userID":   session.Get("userID"),
	})
}

func postCreatePageHandler(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("userID") == nil {
		c.Redirect(http.StatusFound, "/")
	}
	c.HTML(http.StatusOK, "post_create.html", gin.H{
		"username": session.Get("username"),
		"userID":   session.Get("userID"),
	})
}

func postUpdatePageHandler(c *gin.Context) {
	session := sessions.Default(c)
	var article Article
	id := c.Param("id")
	db.First(&article, id)

	if session.Get("userID") != article.UserID {
		c.Redirect(http.StatusFound, "/")
	}

	c.HTML(http.StatusOK, "post_edit.html", gin.H{
		"post":     article,
		"username": session.Get("username"),
		"userID":   session.Get("userID"),
	})
}

func userSignUpPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "user_signup.html", gin.H{})
}

func userSignInPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "user_signin.html", gin.H{})
}

func postCreateHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID").(uint)
	title := c.PostForm("title")
	content := c.PostForm("content")
	db.Create(&Article{Title: title, Content: content, UserID: userID})
	c.Redirect(http.StatusFound, "/post")
}

func postUpdateHandler(c *gin.Context) {
	id := c.Param("id")
	title := c.PostForm("title")
	content := c.PostForm("content")
	db.Where(id).Updates(&Article{Title: title, Content: content})
	c.Redirect(http.StatusFound, "/post")
}

func postDeleteHandler(c *gin.Context) {
	var article Article
	id := c.Param("id")
	db.First(&article, id)
	session := sessions.Default(c)
	if session.Get("userID") != article.UserID {
		c.Redirect(http.StatusFound, "/")
		return
	}
	db.Delete(&article)
	c.Redirect(http.StatusFound, "/post")
}

func userSignUpHandler(c *gin.Context) {
	var existUser User
	username := c.PostForm("username")
	password := c.PostForm("password")
	confirm := c.PostForm("confirm")

	if db.First(&existUser, "username", username).Error != nil {
		if password == confirm {
			hash, _ := HashPassword(password)
			db.Create(&User{Username: username, Password: hash})
			c.Redirect(http.StatusFound, "/user/signin")
		} else {
			c.HTML(http.StatusBadRequest, "user_signup_error.html", gin.H{})
		}
	} else {
		c.HTML(http.StatusBadRequest, "user_signup_error.html", gin.H{})
	}
}

func userSignInHandler(c *gin.Context) {
	var existUser User
	username := c.PostForm("username")
	password := c.PostForm("password")

	db.First(&existUser, "username", username)
	if CheckPasswordHash(password, existUser.Password) {
		session := sessions.Default(c)
		session.Set("username", existUser.Username)
		session.Set("userID", existUser.ID)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}
		fmt.Println(session)
		c.Redirect(http.StatusFound, "/")
	} else {
		c.HTML(http.StatusBadRequest, "user_signin_error.html", gin.H{})
	}
}

func userSignOutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("username", nil)
	session.Set("userID", nil)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.Redirect(http.StatusFound, "/")
}

func filePageHandler(c *gin.Context) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "file_home.html", gin.H{
		"username": session.Get("username"),
		"userID":   session.Get("userID"),
	})
}

func fileUploadHandler(c *gin.Context) {
	file, _ := c.FormFile("file")
	fmt.Println(file.Filename)
	c.SaveUploadedFile(file, filepath.Join("./upload", file.Filename))
	c.Redirect(http.StatusFound, "/file")
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
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.StaticFS("/upload", fs)
	r.LoadHTMLGlob("pages/*")

	private := r.Group("")
	private.Use(AuthRequired)
	{
		r.GET("/post/create", postCreatePageHandler)
		r.POST("/post/create", postCreateHandler)
		r.GET("/post/:id/edit", postUpdatePageHandler)
		r.POST("/post/:id/edit", postUpdateHandler)
		r.GET("/post/:id/delete", postDeleteHandler)
		r.GET("/user/signout", userSignOutHandler)
	}

	r.GET("/", indexPageHandler)
	r.GET("/post", postListPageHandler)
	r.GET("/post/:id", postDetailPageHandler)
	r.GET("/user/signup", userSignUpPageHandler)
	r.POST("/user/signup", userSignUpHandler)
	r.GET("/user/signin", userSignInPageHandler)
	r.POST("/user/signin", userSignInHandler)
	r.GET("/file", filePageHandler)
	r.POST("/file", fileUploadHandler)

	return r
}
