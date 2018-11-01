package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ArticlePusSys/models"
	"path"
	"time"
	"math"
	"github.com/gomodule/redigo/redis"
	"encoding/gob"
	"bytes"
)

type ArticleControllers struct {
	beego.Controller
}

var UserName interface{}

/** 显示文章首页 */
func (this *ArticleControllers) ShowArticleHome() {

	o := orm.NewOrm()
	qs := o.QueryTable("Article")
	var articles []models.Article

	pageIndex,err :=this.GetInt("pageIndex")

	if err != nil {
		pageIndex = 1
	}

	// 获取传递过来选中的文章类型
	typeName := this.GetString("select")

	beego.Info("type name",typeName)

	var count int64
	pageSize := 2
	// 计算获取内容的起始位置
	start := (pageIndex -1)*pageSize

	// 这里是获取内容总数
	if typeName == ""{
		count,_ = qs.Count()
	} else{
		count ,_ = qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()
	}

	pageCount := math.Ceil(float64(count)/float64(pageSize))

	// 文章类型
	var types []models.ArticleType
	conn,err :=redis.Dial("tcp",":6379")
	defer conn.Close()

	buffer,err :=redis.Bytes(conn.Do("get","types"))
	if err != nil {
		beego.Info("获取Redis数据错误")
	}

	dec := gob.NewDecoder(bytes.NewBuffer(buffer))
	err = dec.Decode(&types)
	if len(types) == 0 {
		beego.Info("从数据库中获取数据")
		o.QueryTable("ArticleType").All(&types)
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		err := enc.Encode(&types)
		_,err = conn.Do("set","types",buffer.Bytes())
		if err != nil {
			beego.Info("操作数据库错误")
			return
		}
	}

	// 获取文章内容数据
	if typeName == ""{
		qs.Limit(pageSize,start).All(&articles)
	}else{
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)
	}

beego.Info("type class ",articles[0].ArticleType,typeName)

	this.Data["articles"] = articles
	this.Data["count"] = count
	this.Data["pageIndex"] = pageIndex
	this.Data["pageCount"] = int(pageCount)
	this.Data["types"] = types
	this.Data["typeName"] = typeName

	this.Data["title"] = "文章列表"
	this.Data["userName"] = UserName

	this.Layout = "layout.html"
	this.TplName = "index.html"
}

/** 显示添加文章 */
func (this *ArticleControllers) ShowAddArticle() {

	o := orm.NewOrm()
	var articleTypes []models.ArticleType

	o.QueryTable("ArticleType").All(&articleTypes)

	this.Data["types"] = articleTypes
	this.Data["title"] = "添加文章"
	this.Data["userName"] = UserName

	this.Layout = "layout.html"
	this.TplName = "add.html"
}

/** 添加文章处理 */
func (this *ArticleControllers) HandleAddArticle() {
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	beego.Info(articleName, content)
	if articleName == "" || content == "" {

		this.Data["errmsg"] = "文章数据添加不完整"
		this.Data["title"] = "添加文章"
		this.Data["userName"] = UserName

		this.Layout = "layout.html"
		this.TplName = "add.html"
		return
	}

	fileName := UploadFile(&this.Controller, "uploadname")

	if fileName == "" {
		this.Data["title"] = "添加文章"
		this.Data["userName"] = UserName
		this.Layout = "layout.html"
		this.TplName = "add.html"
		return
	}

	// 封装文章数据
	o := orm.NewOrm()
	var article models.Article
	article.ArtiName = articleName
	article.Acontent = content

	if fileName != "NoImg" && fileName != ""{

	article.Aimg = fileName
	}

	// 获取选择的文章类型
	typeName := this.GetString("select")
	var artycleType models.ArticleType
	artycleType.TypeName = typeName
	// 根据文章名称查询数据
	o.Read(&artycleType, "TypeName")
	// 把对应的类型赋值
	article.ArticleType = &artycleType

	o.Insert(&article)

	this.Redirect("/article/home", 302)
}

/** 显示文章详情 */
func (this *ArticleControllers) ShowArticleDetail() {

	articleId,err := this.GetInt("articleId")
	if err != nil {
		beego.Error("获取文章数据失败")
		return
	}
	o := orm.NewOrm()
	var article models.Article
	article.Id = articleId
	// 关联文章类型表查询，过滤重复的数据
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("id",articleId).One(&article)

	article.Acount +=1
	o.Update(&article)

	// 多表插入
	m2m := o.QueryM2M(&article,"Users")
	var user models.User
	user.Name = UserName.(string)
	// 先获取用户信息
	o.Read(&user,"Name")
	// 插入用户关联信息
	m2m.Add(&user)

	// 关联查询用户信息，Distinct()方法过滤重复的用户数据
	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id",articleId).Distinct().All(&users)
	//o.QueryTable("User").Filter("Articles__Article__Id",articleId).Distinct().All(&users)

beego.Info(article.ArticleType)
	this.Data["users"] = users
	this.Data["article"] = article
	this.Data["title"] = "文章详情"
	this.Data["userName"] = UserName

	this.Layout = "layout.html"
	this.TplName = "content.html"
}

/** 显示文章更新界面 */
func (this *ArticleControllers) ShowUpdateArticle() {

	articleId, err := this.GetInt("articleId")
	if err != nil {
		beego.Info("文章请求错误")
		return
	}

	o := orm.NewOrm()
	var article models.Article
	article.Id = articleId

	o.Read(&article)
	this.Data["article"] = article

	this.Data["title"] = "编辑文章"
	this.Data["userName"] = UserName

	this.Layout = "layout.html"
	this.TplName = "update.html"
}

/** 处理文章更新 */
func (this *ArticleControllers) HandleUpdateArticle() {
	articleId, err := this.GetInt("articleId")

	articleName := this.GetString("articleName")
	content := this.GetString("content")

	fileName := UploadFile(&this.Controller, "uploadname")

	if err != nil || articleName == "" || content == "" || fileName == "" {
		beego.Error("请求数据错误")
		return
	}

	// 封装数据，更新数据库
	o := orm.NewOrm()
	var article models.Article
	article.Id = articleId

	err = o.Read(&article)
	if err != nil {
		beego.Error("更新的文章不存在")
		return
	}
	article.ArtiName = articleName
	article.Acontent = content

	if fileName != "NoImg" && fileName != "" {
	article.Aimg = fileName
	}


	o.Update(&article)

	this.Redirect("/article/home", 302)
}

/** 上传文件 */
func UploadFile(this *beego.Controller, filePath string) string {

	file, head, err := this.GetFile(filePath)
	if head.Filename == ""{
		return "NoImg"
	}
	if err != nil {
		this.Data["errmsg"] = "文件上传失败"
		return ""
	}
	defer file.Close()

	if head.Size > 5000000 {
		this.Data["errmsg"] = "图片文件太大，请重新上传"
		return ""
	}

	// 判断后缀
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		this.Data["errmsg"] = "图片文件格式错误，请重新上传"
		return ""
	}

	fileName := time.Now().Format("2006-01-02-15.04.05")+ext
	this.SaveToFile(filePath, "./static/img"+fileName)

	return "/static/img" + fileName

}

/** 文章删除处理 */
func (this *ArticleControllers) DeleteArticle() {
	articleId,err := this.GetInt("articleId")

	if err != nil {
		beego.Error("文章删除失败")
		return
	}
	o := orm.NewOrm()

	var article models.Article
	article.Id = articleId

	o.Delete(&article)



	this.Redirect("/article/home", 302)
}

/** 显示添加分类 */
func (this *ArticleControllers) ShowAddType() {

	o := orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)

	this.Data["types"] = types
	this.Data["title"] = "添加分类"
	this.Data["userName"] = UserName
	this.Layout = "layout.html"
	this.TplName = "addType.html"
}

/** 处理添加分类 */
func (this *ArticleControllers) HandleAddType() {

	typeName := this.GetString("typeName")
	if typeName == "" {
		this.Data["errmsg"] = "请输入内容"
		return
	}

	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Insert(&articleType)

	this.Redirect("/article/addType", 302)
}

/** 删除类型 */
func (this *ArticleControllers) DeleteType() {
	id, err := this.GetInt("typeId")
	if err != nil {
		beego.Info("删除文章类型失败")
		return
	}

	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id = id
	o.Delete(&articleType)

	this.Redirect("/article/addType", 302)
}
