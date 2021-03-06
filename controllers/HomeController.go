package controllers

import (
	"strings"

	"sdrms/enums"
	"sdrms/models"
	"sdrms/utils"
)

type HomeController struct {
	BaseController
}

func (c *HomeController) Index() {
	//判断是否登录
	c.checkLogin()
	// c.GetSession()
	c.setTpl()
}
func (c *HomeController) Page404() {
	c.setTpl()
}
func (c *HomeController) Error() {
	c.Data["error"] = c.GetString(":error")
	c.setTpl("home/error.html", "shared/layout_pullbox.html")
}

// func (c *HomeController) Login() {

// 	c.LayoutSections = make(map[string]string)
// 	c.LayoutSections["headcssjs"] = "home/login_headcssjs.html"
// 	c.LayoutSections["footerjs"] = "home/login_footerjs.html"
// 	c.setTpl("home/login.html", "shared/layout_base.html")
// }

//new login function
func (c *HomeController) Login() {
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "home/admin_login_headcssjs.html"
	c.LayoutSections["footerjs"] = "home/admin_login_footerjs.html"
	c.setTpl("home/admin.html", "shared/layout_base_login.html")

}

func (c *HomeController) DoLogin() {
	username := strings.TrimSpace(c.GetString("UserName"))
	userpwd := strings.TrimSpace(c.GetString("UserPwd"))
	if len(username) == 0 || len(userpwd) == 0 {
		c.jsonResult(enums.JRCodeFailed, "用户名和密码不正确", "")
	}
	userpwd = utils.String2md5(userpwd)
	user, err := models.BackendUserOneByUserName(username, userpwd)
	if user != nil && err == nil {
		if user.Status == enums.Disabled {
			c.jsonResult(enums.JRCodeFailed, "用户被禁用，请联系管理员", "")
		}
		//保存用户信息到session
		c.setBackendUser2Session(user.Id)
		// println(c.setBackendUser2Session(user.Id))
		//获取用户信息
		c.jsonResult(enums.JRCodeSucc, "登录成功", "")
	} else {
		c.jsonResult(enums.JRCodeFailed, "用户名或者密码错误", "")
	}
}
func (c *HomeController) Logout() {
	user := models.BackendUser{}
	c.SetSession("backenduser", user)
	c.pageLogin()
}
