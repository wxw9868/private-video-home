package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func NoRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
		c.Abort()
	}
}

func InitSession() gin.HandlerFunc {
	store := cookie.NewStore([]byte("secret_20240109"))
	return sessions.Sessions("stock_session", store)
}

func AuthSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		userid, ok := session.Get("userID").(uint)
		fmt.Println("user: ", userid, ok)
		if !ok {
			fmt.Println("user out: ", userid, ok)
			c.Redirect(http.StatusMovedPermanently, "/login")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
