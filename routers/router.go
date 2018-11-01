package routers

import (
	"ArticlePusSys/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	/*
		过滤器，对应有5种状态
		BeforeStatic	静态地址之前
		BeforeRouter	寻找路由之前
		BeforeExec	找到路由之后，开始执行相应的	Controller	之前
		AfterExec	执行完	Controller	逻辑之后执行的过滤器
		FinishRouter	执行完逻辑之后执行的过滤器
	*/
	beego.InsertFilter("/article/*", beego.BeforeExec, Firfter)

	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.UserControllers{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/register", &controllers.UserControllers{}, "get:ShowRegister;post:HandleRegister")
	beego.Router("/logout", &controllers.UserControllers{}, "get:Logout")

	// 后台主页，显示文章列表页
	beego.Router("/article/home", &controllers.ArticleControllers{}, "get:ShowArticleHome")
	// 添加文章页
	beego.Router("/article/addArticle", &controllers.ArticleControllers{}, "get:ShowAddArticle;post:HandleAddArticle")
	// 文章详情
	beego.Router("/article/detailArticle", &controllers.ArticleControllers{}, "get:ShowArticleDetail")
	// 编辑文章
	beego.Router("/article/updateArticle", &controllers.ArticleControllers{}, "get:ShowUpdateArticle;post:HandleUpdateArticle")
	// 删除文章
	beego.Router("/article/deleteArticle", &controllers.ArticleControllers{}, "get:DeleteArticle")
	// 添加文章类型
	beego.Router("/article/addType", &controllers.ArticleControllers{}, "get:ShowAddType;post:HandleAddType")
	// 删除文章类型
	beego.Router("/article/deleteType", &controllers.ArticleControllers{}, "get:DeleteType")
}

var Firfter = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")

	if userName == nil {
		ctx.Redirect(302, "/login")
		return
	}
	controllers.UserName = userName
}
