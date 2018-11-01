package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ArticlePusSys/models"
)

type UserControllers struct {
	beego.Controller
}

/** 显示登录界面 */
func (this *UserControllers) ShowLogin() {
	// 获取记住的用户名
	userName := this.Ctx.GetCookie("userName")

	this.Data["userName"] = userName
	this.Data["checked"] = "checked"
	// 没有获取到数据，返回空的过去
	if userName == "" {
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	}

	this.TplName = "login.html"
}

/** 用户登录操作 */
func (this *UserControllers) HandleLogin() {
	userName := this.GetString("userName")
	password := this.GetString("password")

	if userName == "" || password == "" {
		this.Data["errmsg"] = "用户名或密码为空"
		this.TplName = "login.html"
		return
	}

	// 获取orm对象
	o := orm.NewOrm()

	// 处理数据
	var user models.User
	user.Name = userName
	user.Password = password

	err := o.Read(&user, "Name")
	// 判断读取用户信息
	if err != nil {
		this.Data["errmsg"] = "用户不存在"
		this.TplName = "login.html"
		return
	}

	// 判断用户密码
	if user.Password != password {
		this.Data["errmsg"] = "用户密码错误"
		this.TplName = "login.html"
		return
	}
	// 判断是否记住用户名
	data := this.GetString("remember")
	if data == "on" {
		this.Ctx.SetCookie("userName", userName, 1000*60*60)
	} else {
		this.Ctx.SetCookie("userName", userName, -1)
	}

	this.SetSession("userName", userName)

	this.Redirect("/article/home", 302)
}

/** 显示注册界面 */
func (this *UserControllers) ShowRegister() {
	this.TplName = "register.html"
}

/** 用户注册操作 */
func (this *UserControllers) HandleRegister() {
	userName := this.GetString("userName")
	password := this.GetString("password")

	if userName == "" || password == "" {
		this.Data["errmsg"] = "用户名或密码为空"
		this.TplName = "register.html"
		return
	}

	o := orm.NewOrm()

	var user models.User
	user.Name = userName
	user.Password = password

	_, err := o.Insert(&user)

	if err != nil {
		beego.Error("用户注册失败")
		this.Data["errmsg"] = "用户注册失败，请重新注册"
		this.TplName = "register.html"
		return
	}

	this.Redirect("/login", 302)
}

/** 用户退出 */
func (this *UserControllers) Logout() {

	this.DelSession("userName")
	this.Redirect("/login", 302)
}
