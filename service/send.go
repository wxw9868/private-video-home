package service

import (
	"github.com/wxw9868/util/mail"
)

type SendService struct{}

// SendMail 发送邮箱验证码
func (ss *SendService) SendMail(emails []string, content, code string) error {
	conf := `{"username":"986845663@qq.com","password":"emtpyouqirhebfij","host":"smtp.qq.com","port":587}`
	m := mail.NewEmail(conf)
	m.To = emails
	m.From = "986845663@qq.com"
	m.Subject = "验证码"
	m.Text = "Text Body is, of course, supported!"
	m.HTML = `
	<head>
		<base target="_blank" />
		<style type="text/css">
			::-webkit-scrollbar {
				display: none;
			}
		</style>
		<style id="cloudAttachStyle" type="text/css">
			#divNeteaseBigAttach,
			#divNeteaseBigAttach_bak {
				display: none;
			}
		</style>
		<style id="blockquoteStyle" type="text/css">
			blockquote {
				display: none;
			}
		</style>
		<style type="text/css">
			body {
				font-size: 14px;
				font-family: arial, verdana, sans-serif;
				line-height: 1.666;
				padding: 0;
				margin: 0;
				overflow: auto;
				white-space: normal;
				word-wrap: break-word;
				min-height: 100px
			}

			td,
			input,
			button,
			select,
			body {
				font-family: Helvetica, 'Microsoft Yahei', verdana
			}

			pre {
				white-space: pre-wrap;
				white-space: -moz-pre-wrap;
				white-space: -pre-wrap;
				white-space: -o-pre-wrap;
				word-wrap: break-word;
				width: 95%
			}

			th,
			td {
				font-family: arial, verdana, sans-serif;
				line-height: 1.666
			}

			img {
				border: 0
			}

			header,
			footer,
			section,
			aside,
			article,
			nav,
			hgroup,
			figure,
			figcaption {
				display: block
			}

			blockquote {
				margin-right: 0px
			}
		</style>
	</head>

	<body tabindex="0" role="listitem">
		<table width="700" border="0" align="center" cellspacing="0" style="width:700px;">
			<tbody>
				<tr>
					<td>
						<div style="width:700px;margin:0 auto;border-bottom:1px solid #ccc;margin-bottom:30px;">
							<table border="0" cellpadding="0" cellspacing="0" width="700" height="39" style="font:12px Tahoma, Arial, 宋体;">
								<tbody>
									<tr>
										<td width="210"></td>
									</tr>
								</tbody>
							</table>
						</div>
						<div style="width:680px;padding:0 10px;margin:0 auto;">
							<div style="line-height:1.5;font-size:14px;margin-bottom:25px;color:#4d4d4d;">
								<strong style="display:block;margin-bottom:15px;">尊敬的用户：<span style="color:#f60;font-size: 16px;"></span>您好！</strong>
								<strong style="display:block;margin-bottom:15px;">
									您正在进行<span style="color: red">` + content + `</span>操作，请在验证码输入框中输入：<span style="color:#f60;font-size: 24px">` + code + `</span>，以完成操作。
								</strong>
							</div>
							<div style="margin-bottom:30px;">
								<small style="display:block;margin-bottom:20px;font-size:12px;">
									<p style="color:#747474;">
										注意：此操作可能会修改您的密码、登录邮箱或绑定手机。如非本人操作，请及时登录并修改密码以保证帐户安全<br>
										（工作人员不会向你索取此验证码，请勿泄漏！）
									</p>
								</small>
							</div>
						</div>
						<div style="width:700px;margin:0 auto;">
							<div style="padding:10px 10px 0;border-top:1px solid #ccc;color:#747474;margin-bottom:20px;line-height:1.3em;font-size:12px;">
								<p>
									此为系统邮件，请勿回复<br>
									请保管好您的邮箱，避免账号被他人盗用
								</p>
								<p> </p>
							</div>
						</div>
					</td>
				</tr>
			</tbody>
		</table>
	</body>`
	if err := m.Send(); err != nil {
		return err
	}
	return nil
}

// SendUrl 发送一个唯一的、有时效性的密码重置链接
func (ss *SendService) SendUrl(emails []string, code string) error {
	conf := `{"username":"986845663@qq.com","password":"emtpyouqirhebfij","host":"smtp.qq.com","port":587}`
	m := mail.NewEmail(conf)
	m.To = emails
	m.From = "986845663@qq.com"
	m.Subject = "重置密码"
	m.Text = "Text Body is, of course, supported!"
	m.HTML = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Reset Your Ledger Password</title>
	<style>
		body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f4f4f4; }
		.container { background-color: #ffffff; max-width: 500px; margin: 40px auto; padding: 20px; box-shadow: 0 0 5px rgba(0,0,0,0.2); }
		.header { background-color: #ffffff; color: #24292e; padding: 0px; text-align: center; }
		.content { padding: 0px; text-align: center; }
		.button { background-color: #0366d6; color: #ffffff; padding: 15px 30px; margin: 0 0; display: inline-block; text-decoration: none; border-radius: 5px; }
		.footer { margin-top: 10px; text-align: center; font-size: 12px; color: #666666; }
	</style>
	</head>
	<body>
	<div class="container">
		<div class="header">
			<h1>Video System</h1>
		</div>
		<div class="content">
			<h2>Reset Your Password</h2>
			<p>We heard that you lost your Ledger password. Sorry about that! But don’t worry! You can use the following button to reset your password:</p>
			<a href=http://127.0.0.1:80/reset-pwd?reset_password_token=` + code + ` class="button">Reset Password</a>
			<p>If you don’t use this link within 3 hours, it will expire. To get a new password reset link, visit: </p>
			<p><a href="http://127.0.0.1:8080/forgot-pwd">http://127.0.0.1:8080/forgot-pwd</a></p>
		</div>
		<div class="footer">
			<p>Thank you for using Ledger!</p>
		</div>
	</div>
	</body>
	</html>`
	if err := m.Send(); err != nil {
		return err
	}
	return nil
}
