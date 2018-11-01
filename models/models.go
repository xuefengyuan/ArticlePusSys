package models

import "github.com/astaxie/beego/orm"
import (
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id       int
	Name     string
	Password string
	Articles []*Article `orm:"reverse(many)"`
}

type Article struct {
	Id int`orm:"pk;auto"`
	ArtiName string `orm:"size(20)"`
	Atime time.Time `orm:"auto_now"`
	Acount int `orm:"default(0);null"`
	Acontent string `orm:"size(500)"`
	Aimg string  `orm:"size(100)"`

	ArticleType *ArticleType `orm:"rel(fk);on_delete(set_null);null"`
	Users []*User `orm:"rel(m2m)"`
}

type ArticleType struct {
	Id int
	TypeName string `orm:"size(20)"`

	Articles []*Article `orm:"reverse(many)"`
}


func init() {
	// ORM 操作数据库
	// 获取连接对象
	orm.RegisterDataBase("default", "mysql", "root:12345678@tcp(192.168.1.5:3306)/test1?charset=utf8")
	// 创建表
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	// 生成表，第一个参数是数据库别名，第二个参数是否是强制更新，设置为true则每次会清空数据，第三个参数显创建表过程
	orm.RunSyncdb("default", false, true)
	//
}
