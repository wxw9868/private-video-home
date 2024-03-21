package api

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LoginApi(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "登录",
	})
}

func DoLoginApi(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	if email != "" && password != "" {
		session := sessions.Default(c)
		session.Set("email", email)
		session.Set("password", password)
		session.Save()

		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
}

func LogoutApi(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.JSON(http.StatusOK, nil)
}
