package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wxw9868/util"
	"github.com/wxw9868/video/model"
)

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
	session.Set("userDesignation", user.Designation)
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

func SessionApi(c *gin.Context) {
	session := sessions.Default(c)
	data := map[string]interface{}{
		"id":          session.Get("userID").(uint),
		"avatar":      session.Get("userAvatar").(string),
		"username":    session.Get("userUsername").(string),
		"nickname":    session.Get("userNickname").(string),
		"email":       session.Get("userEmail").(string),
		"mobile":      session.Get("userMobile").(string),
		"designation": session.Get("userDesignation").(string),
	}
	c.JSON(http.StatusOK, util.Success("获取成功", data))
}

func UserInfoApi(c *gin.Context) {
	user, err := us.Info(GetUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("获取成功", user))
}

type ChangePassword struct {
	OldPassword     string `form:"old_password" json:"old_password" binding:"required"`
	NewPassword     string `form:"new_password" json:"new_password" binding:"required"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" binding:"required,eqcsfield=NewPassword"`
}

func ChangePasswordApi(c *gin.Context) {
	var bind ChangePassword
	if err := c.ShouldBind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	userID := GetUserID(c)
	if err := us.ChangePassword(userID, bind.OldPassword, bind.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("修改密码成功", nil))
}

type ForgotPassword struct {
	ResetPasswordToken string `form:"reset_password_token" json:"reset_password_token" binding:"required"`
	Password           string `form:"password" json:"password" binding:"required"`
	ConfirmPassword    string `form:"confirm_password" json:"confirm_password" binding:"required,eqcsfield=Password"`
}

func ForgotPasswordApi(c *gin.Context) {
	var bind ForgotPassword
	if err := c.ShouldBind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	session := sessions.Default(c)
	email := session.Get(bind.ResetPasswordToken).(string)
	if email == "" {
		c.JSON(http.StatusBadRequest, util.Fail("密码重置链接已失效，请重新获取"))
		return
	}

	if err := us.ForgotPassword(email, bind.Password); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("修改密码成功", nil))
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

	user := model.User{
		Username:    bind.Username,
		Nickname:    bind.Nickname,
		Mobile:      bind.Mobile,
		Email:       bind.Email,
		Designation: bind.Designation,
		Address:     bind.Address,
		Note:        bind.Note,
	}
	err := us.Updates(GetUserID(c), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("更新成功", nil))
}

func UserUploadAvatarApi(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(fmt.Sprintf("get form err: %s", err.Error())))
		return
	}

	filename := filepath.Base(file.Filename)
	avatarDir := "./assets/image/avatar/" + filename
	if err := c.SaveUploadedFile(file, avatarDir); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(fmt.Sprintf("upload file err: %s", err.Error())))
		return
	}

	session := sessions.Default(c)
	if err := us.Update(session.Get("userID").(uint), "avatar", avatarDir); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	if session.Get("userAvatar").(string) != "./assets/image/avatar/avatar.png" {
		os.Remove(session.Get("userAvatar").(string))
	}
	session.Set("userAvatar", avatarDir)
	session.Save()

	c.JSON(http.StatusOK, util.Success("更换成功", avatarDir))
}

func UserCollectApi(c *gin.Context) {
	data, err := us.CollectList(GetUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("获取成功", data))
}

func UserBrowseApi(c *gin.Context) {
	data, err := us.BrowseList(GetUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("获取成功", data))
}
