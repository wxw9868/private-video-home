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

type Register struct {
	Username       string `form:"username" json:"username" binding:"required"`
	Email          string `form:"email" json:"email" binding:"required,email"`
	Password       string `form:"password" json:"password" binding:"required"`
	RepeatPassword string `form:"repeat_password" json:"repeat_password" binding:"required,eqcsfield=Password"`
}

// RegisterApi godoc
//
//	@Summary	用户注册
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		data	body		Register	true	"用户注册信息"
//	@Success	200		{object}	Success
//	@Failure	400		{object}	Fail
//	@Failure	404		{object}	NotFound
//	@Failure	500		{object}	ServerError
//	@Router		/user/register [post]
func RegisterApi(c *gin.Context) {
	var bind Register
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := util.VerifyPassword(bind.Password); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	err := userService.Register(bind.Username, bind.Email, bind.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("注册成功", nil))
}

type Login struct {
	Email    string `form:"email" json:"email" binding:"required,email" example:"wxw9868@163.com"`
	Password string `form:"password" json:"password" binding:"required" example:"123456"`
}

// LoginApi godoc
//
//	@Summary		用户登录
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			data	body		Login	true	"用户登录信息"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/user/login [post]
func LoginApi(c *gin.Context) {
	var bind Login
	if err := c.ShouldBind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	user, err := userService.Login(bind.Email, bind.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("user_avatar", user.Avatar)
	session.Set("user_username", user.Username)
	session.Set("user_nickname", user.Nickname)
	session.Set("user_email", user.Email)
	session.Set("user_mobile", user.Mobile)
	session.Set("user_designation", user.Designation)
	if err = session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail("登录失败"))
		return
	}

	c.JSON(http.StatusOK, util.Success("登录成功", nil))
}

// LogoutApi godoc
//
//	@Summary		用户登出
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Success
//	@Failure		404	{object}	NotFound
//	@Failure		500	{object}	ServerError
//	@Router			/user/session [get]
func LogoutApi(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("登出成功", nil))
}

// SessionApi godoc
//
//	@Summary		用户信息
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Success
//	@Failure		404	{object}	NotFound
//	@Router			/user/session [get]
func SessionApi(c *gin.Context) {
	session := sessions.Default(c)
	data := map[string]interface{}{
		"id":          session.Get("user_id").(uint),
		"avatar":      session.Get("user_avatar").(string),
		"username":    session.Get("user_username").(string),
		"nickname":    session.Get("user_nickname").(string),
		"email":       session.Get("user_email").(string),
		"mobile":      session.Get("user_mobile").(string),
		"designation": session.Get("user_designation").(string),
	}
	c.JSON(http.StatusOK, util.Success("获取成功", data))
}

// GetUserInfoApi godoc
//
//	@Summary		获取用户信息
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Success
//	@Failure		404	{object}	NotFound
//	@Failure		500	{object}	ServerError
//	@Router			/user/getUserInfo [get]
func GetUserInfoApi(c *gin.Context) {
	user, err := userService.Info(GetUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("SUCCESS", user))
}

type ChangePassword struct {
	OldPassword     string `form:"old_password" json:"old_password" binding:"required"`
	NewPassword     string `form:"new_password" json:"new_password" binding:"required"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" binding:"required,eqcsfield=NewPassword"`
}

// ChangePasswordApi godoc
//
//	@Summary		修改密码
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			data	body		ChangePassword	true	"修改密码信息"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/user/changePassword [post]
func ChangePasswordApi(c *gin.Context) {
	var bind ChangePassword
	if err := c.ShouldBind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	userID := GetUserID(c)
	if err := userService.ChangePassword(userID, bind.OldPassword, bind.NewPassword); err != nil {
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

// ForgotPasswordApi godoc
//
//	@Summary		忘记密码
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			data	body		ForgotPassword	true	"忘记密码信息"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/user/forgotPassword [post]
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

	if err := userService.ForgotPassword(email, bind.Password); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("修改密码成功", nil))
}

type UpdateUserInfo struct {
	Nickname    string `form:"nickname" json:"nickname" binding:"required"`
	Username    string `form:"username" json:"username" binding:"required"`
	Email       string `form:"email" json:"email" binding:"required,email"`
	Mobile      string `form:"mobile" json:"mobile" binding:"required"`
	Designation string `form:"designation" json:"designation"`
	Address     string `form:"address" json:"address"`
	Note        string `form:"note" json:"note"`
}

// UpdateUserInfoApi godoc
//
//	@Summary		修改用户信息
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			data	body		UserUpdate	true	"修改用户信息"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/user/updateUserInfo [post]
func UpdateUserInfoApi(c *gin.Context) {
	var bind UpdateUserInfo
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
	err := userService.Updates(GetUserID(c), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("更新成功", nil))
}

// ChangeUserAvatarApi godoc
//
//	@Summary		更换用户头像
//	@Description	上传头像文件
//	@Tags			user
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			avatar	formData	file	true	"头像文件"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/user/changeUserAvatar [post]
func ChangeUserAvatarApi(c *gin.Context) {
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
	if err := userService.Update(session.Get("user_id").(uint), "avatar", avatarDir); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	if session.Get("user_avatar").(string) != "./assets/image/avatar/avatar.png" {
		os.Remove(session.Get("user_avatar").(string))
	}
	session.Set("user_avatar", avatarDir)
	session.Save()

	c.JSON(http.StatusOK, util.Success("更换成功", avatarDir))
}

// GetUserFavoriteListApi godoc
//
//	@Summary		获取用户收藏记录
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Success
//	@Failure		404	{object}	NotFound
//	@Failure		500	{object}	ServerError
//	@Router			/user/getUserFavoriteList [get]
func GetUserFavoriteListApi(c *gin.Context) {
	data, err := userService.FavoriteList(GetUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("获取用户收藏记录", data))
}

// GetUserBrowseListApi godoc
//
//	@Summary		获取用户浏览记录
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Success
//	@Failure		404	{object}	NotFound
//	@Failure		500	{object}	ServerError
//	@Router			/user/getUserBrowseHistory [get]
func GetUserBrowseListApi(c *gin.Context) {
	data, err := userService.BrowseList(GetUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("获取用户浏览记录", data))
}
