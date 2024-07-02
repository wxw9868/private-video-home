package api

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wxw9868/util"
	"github.com/wxw9868/video/model"
)

func Register(c *gin.Context) {
	c.HTML(http.StatusOK, "sign-up.html", gin.H{
		"title": "注册",
	})
}

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "sign-in.html", gin.H{
		"title": "登录",
	})
}

func Account(c *gin.Context) {
	c.HTML(http.StatusOK, "account.html", gin.H{
		"title": "账号设置",
	})
}

type RegisterReq struct {
	Username       string `form:"username" json:"username" binding:"required"`
	Email          string `form:"email" json:"email" binding:"required,email"`
	Password       string `form:"password" json:"password" binding:"required"`
	RepeatPassword string `form:"repeat_password" json:"repeat_password" binding:"required,eqcsfield=Password"`
}

func RegisterApi(c *gin.Context) {
	var bind RegisterReq
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := util.VerifyPassword(bind.Password); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	err := us.Register(bind.Username, bind.Email, bind.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("注册成功", nil))
}

type LoginReq struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

func LoginApi(c *gin.Context) {
	var bind LoginReq
	if err := c.ShouldBind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	user, err := us.Login(bind.Email, bind.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	session := sessions.Default(c)
	session.Set("userID", user.ID)
	session.Set("userAvatar", user.Avatar)
	session.Set("userUsername", user.Username)
	session.Set("userNickname", user.Nickname)
	session.Set("userEmail", user.Email)
	session.Set("userMobile", user.Mobile)
	if err = session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail("登录失败"))
		return
	}

	c.JSON(http.StatusOK, util.Success("登录成功", nil))
}

func LogoutApi(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("登出成功", nil))
}

func GetSession(c *gin.Context) {
	session := sessions.Default(c)
	data := map[string]interface{}{
		"id":       session.Get("userID").(uint),
		"avatar":   session.Get("userAvatar").(string),
		"username": session.Get("userUsername").(string),
		"nickname": session.Get("userNickname").(string),
		"email":    session.Get("userEmail").(string),
		"mobile":   session.Get("userMobile").(string),
	}
	c.JSON(http.StatusOK, util.Success("获取成功", data))
}

type UserUpdate struct {
	Nickname    string `form:"nickname" json:"nickname" binding:"required"`
	Username    string `form:"username" json:"username" binding:"required"`
	Email       string `form:"email" json:"email" binding:"required,email"`
	Mobile      string `form:"mobile" json:"mobile" binding:"required"`
	Designation string `form:"designation" json:"designation"`
	Address     string `form:"address" json:"address"`
	Note        string `form:"note" json:"note"`
}

func UserUpdateApi(c *gin.Context) {
	var bind UserUpdate
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	session := sessions.Default(c)
	user := model.User{
		Username:    bind.Username,
		Nickname:    bind.Nickname,
		Mobile:      bind.Mobile,
		Email:       bind.Email,
		Designation: bind.Designation,
		Address:     bind.Address,
		Note:        bind.Note,
	}
	err := us.Updates(session.Get("userID").(uint), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

}
