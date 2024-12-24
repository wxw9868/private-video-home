package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"github.com/wxw9868/util"
	"github.com/wxw9868/video/model"
	"github.com/wxw9868/video/model/request"
	"github.com/wxw9868/video/utils"
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
//	@Summary	用户登录
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		data	body		Login	true	"用户登录信息"
//	@Success	200		{object}	Success
//	@Failure	400		{object}	Fail
//	@Failure	404		{object}	NotFound
//	@Failure	500		{object}	ServerError
//	@Router		/user/login [post]
func LoginApi(c *gin.Context) {
	var bind Login
	if err := c.ShouldBind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	user, err := userService.Login(bind.Email, bind.Password)

	ua := user_agent.New(c.Request.UserAgent())
	browser, _ := ua.Browser()
	userLoginLog := model.UserLoginLog{
		LoginName:     bind.Email,
		LoginIpaddr:   c.ClientIP(),
		LoginLocation: utils.GetCityByIp(c.ClientIP()),
		Browser:       browser,
		Os:            ua.OS(),
		LoginTime:     time.Now(),
	}

	if err != nil {
		userLoginLog.UserID = 0
		userLoginLog.Status = 1
		userService.CreateUserLoginLog(userLoginLog)
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	userLoginLog.UserID = user.ID
	userLoginLog.Status = 0
	userService.CreateUserLoginLog(userLoginLog)

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
//	@Summary	用户登出
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	Success
//	@Failure	404	{object}	NotFound
//	@Failure	500	{object}	ServerError
//	@Router		/user/logout [get]
func LogoutApi(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("登出成功", nil))
}

// GetSessionApi godoc
//
//	@Summary	获取用户信息
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	Success
//	@Failure	404	{object}	NotFound
//	@Router		/user/getSession [get]
func GetSessionApi(c *gin.Context) {
	c.JSON(http.StatusOK, util.Success("获取成功", GetUser(c)))
}

// GetUserInfoApi godoc
//
//	@Summary	获取用户信息
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	Success
//	@Failure	404	{object}	NotFound
//	@Failure	500	{object}	ServerError
//	@Router		/user/getUserInfo [get]
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
//	@Summary	修改密码
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		data	body		ChangePassword	true	"修改密码信息"
//	@Success	200		{object}	Success
//	@Failure	400		{object}	Fail
//	@Failure	404		{object}	NotFound
//	@Failure	500		{object}	ServerError
//	@Router		/user/changePassword [post]
func ChangePasswordApi(c *gin.Context) {
	var bind ChangePassword
	if err := c.ShouldBind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := userService.ChangePassword(GetUserID(c), bind.OldPassword, bind.NewPassword); err != nil {
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
//	@Summary	忘记密码
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		data	body		ForgotPassword	true	"忘记密码信息"
//	@Success	200		{object}	Success
//	@Failure	400		{object}	Fail
//	@Failure	404		{object}	NotFound
//	@Failure	500		{object}	ServerError
//	@Router		/user/forgotPassword [post]
func ForgotPasswordApi(c *gin.Context) {
	var bind ForgotPassword
	if err := c.ShouldBind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	session := sessions.Default(c)
	email, ok := session.Get(bind.ResetPasswordToken).(string)
	if email == "" || !ok {
		c.JSON(http.StatusBadRequest, util.Fail("密码重置链接已失效，请重新获取"))
		return
	}

	if err := userService.ForgotPassword(email, bind.Password); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("修改密码成功", nil))
}

// UpdateUserInfoApi godoc
//
//	@Summary	修改用户信息
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		data	body		request.UpdateUser	true	"修改用户信息"
//	@Success	200		{object}	Success
//	@Failure	400		{object}	Fail
//	@Failure	404		{object}	NotFound
//	@Failure	500		{object}	ServerError
//	@Router		/user/updateUserInfo [post]
func UpdateUserInfoApi(c *gin.Context) {
	var bind request.UpdateUser
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
		Intro:       bind.Intro,
	}
	user.ID = GetUserID(c)
	err := userService.Updates(user)
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

	avatarPath := "assets/image/avatar/" + filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, avatarPath); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(fmt.Sprintf("upload file err: %s", err.Error())))
		return
	}

	user := GetUser(c)
	if err := userService.Update(user.ID, "avatar", avatarPath); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	if user.Avatar != "./assets/image/avatar/avatar.png" {
		os.Remove(user.Avatar)
	}

	session := sessions.Default(c)
	session.Set("user_avatar", avatarPath)
	session.Save()

	c.JSON(http.StatusOK, util.Success("更换成功", avatarPath))
}

// GetUserVideoCollectListApi godoc
//
//	@Summary	获取用户收藏记录
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	Success
//	@Failure	404	{object}	NotFound
//	@Failure	500	{object}	ServerError
//	@Router		/user/getUserVideoCollectList [get]
func GetUserVideoCollectListApi(c *gin.Context) {
	data, err := userService.CollectList(GetUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("获取用户收藏记录", data))
}

// GetUserVideoPageViewsListApi godoc
//
//	@Summary	获取用户视频浏览记录
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	Success
//	@Failure	404	{object}	NotFound
//	@Failure	500	{object}	ServerError
//	@Router		/user/getUserVideoPageViewsList [get]
func GetUserVideoPageViewsListApi(c *gin.Context) {
	data, err := userService.VideoPageViewsList(GetUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("获取用户浏览记录", data))
}

// GetUserLoginLogListApi godoc
//
//	@Summary	获取用户登陆记录
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	Success
//	@Failure	404	{object}	NotFound
//	@Failure	500	{object}	ServerError
//	@Router		/user/getUserLoginLogListApi [get]
func GetUserLoginLogListApi(c *gin.Context) {
	data, err := userService.UserLoginLogList(GetUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("获取用户登陆记录", data))
}
