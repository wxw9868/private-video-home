package api

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wxw9868/util"
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
		res, err := us.Login(email, password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.Fail("登录失败"))
			return
		}
		fmt.Printf("%+v\n", res)

		session := sessions.Default(c)
		session.Set("userid", res.ID)
		session.Set("email", res.Email)
		session.Set("password", res.Password)
		if err = session.Save(); err != nil {
			fmt.Println("err: ", err)
			c.JSON(http.StatusInternalServerError, util.Fail("登录失败"))
			return
		}

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
