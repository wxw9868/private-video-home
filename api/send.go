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
	if err := ss.SendMail([]string{email}, "注册账号", code); err != nil {
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

	reset_password_token := uuid.New().String()
	if err := ss.SendUrl([]string{email}, reset_password_token); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail("邮件发送失败！"))
		return
	}

	session := sessions.Default(c)
	session.Set("reset_password_token_"+email, reset_password_token)
	session.Set(reset_password_token, email)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail("发送失败！"))
		return
	}
	c.JSON(http.StatusOK, util.Msg(true, 60002, "发送成功。", nil))
}

type Captcha struct {
	CaptchaType string `from:"captcha_type" json:"captcha_type" binding:"oneof=audio string math chinese"`
}

func CaptchaApi(c *gin.Context) {
	var bind Captcha
	if err := c.ShouldBind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}
	data := captcha.GetCaptcha(bind.CaptchaType)
	c.JSON(http.StatusOK, util.Msg(true, 1, "验证码获取成功！", data))
}
