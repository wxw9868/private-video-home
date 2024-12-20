package api

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wxw9868/util"
	"github.com/wxw9868/util/captcha"
)

type SendMail struct {
	Email string `form:"email" json:"email" binding:"required,email"`
}

// SendMailApi godoc
//
//	@Summary		发送邮件
//	@Description	通过邮件发送验证码
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			data	body		SendMail	true	"发送邮件"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/verify/sendMail [post]
func SendMailApi(c *gin.Context) {
	var bind SendMail
	if err := c.ShouldBind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}
	email := bind.Email
	if !util.VerifyEmail(email) {
		c.JSON(http.StatusBadRequest, util.Fail("邮箱地址格式错误！"))
		return
	}
	code := util.GenerateCode(6)
	if err := sendService.SendMail([]string{email}, "注册账号", code); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail("邮件发送失败！"))
		return
	}
	session := sessions.Default(c)
	session.Set(bind.Email, code)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail("验证码发送失败！"))
		return
	}
	c.JSON(http.StatusOK, util.Msg(true, 60001, "验证码发送成功。", nil))
}

type SendUrl struct {
	Email string `form:"email" json:"email" binding:"required,email"`
}

// SendUrlApi godoc
//
//	@Summary		发送URL
//	@Description	通过邮件发送重置密码连接
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			data	body		SendMail	true	"SendUrl"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/verify/sendUrl [post]
func SendUrlApi(c *gin.Context) {
	var bind SendUrl
	if err := c.ShouldBind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}
	email := bind.Email
	if !util.VerifyEmail(email) {
		c.JSON(http.StatusBadRequest, util.Fail("邮箱地址格式错误！"))
		return
	}

	resetPasswordToken := uuid.New().String()
	if err := sendService.SendUrl([]string{email}, resetPasswordToken); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail("邮件发送失败！"))
		return
	}

	session := sessions.Default(c)
	session.Set("reset_password_token_"+email, resetPasswordToken)
	session.Set(resetPasswordToken, email)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail("发送失败！"))
		return
	}
	c.JSON(http.StatusOK, util.Msg(true, 60002, "发送成功。", nil))
}

type Captcha struct {
	CaptchaType string `from:"captcha_type" json:"captcha_type" binding:"oneof=audio string math chinese"`
}

// CaptchaApi godoc
//
//	@Summary		验证码
//	@Description	验证码
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			data	body		Captcha	true	"验证码"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/verify/captcha [post]
func CaptchaApi(c *gin.Context) {
	var bind Captcha
	if err := c.ShouldBind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}
	data := captcha.GetCaptcha(bind.CaptchaType)
	c.JSON(http.StatusOK, util.Msg(true, 1, "验证码获取成功！", data))
}
