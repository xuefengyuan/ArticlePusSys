package main

import (
	_ "ArticlePusSys/routers"
	"github.com/astaxie/beego"
	_ "ArticlePusSys/models"
)

func main() {
	beego.AddFuncMap("prepage",ShowPrePage)
	beego.AddFuncMap("nextpage",ShowNextPage)
	beego.Run()
}

/**
	上一页
	前端调用，只有一个参数时，参数在前面 方法在后面 {{.参数名 | 方法名}}
	{{.pageIndex | prepage}}
 */
func ShowPrePage(pageIndex int)int{
	if pageIndex == 1{
		return pageIndex
	}
	return pageIndex -1
}

/**
  下一页
  前端调用，有多个参数时，方法在前面，参数在后面 {{方法名 .参数名1 .参数名2 ...}}
 {{ nextpage .pageIndex .pageCount}}
 */
func ShowNextPage(pageIndex,pageCount int)int{
	if pageIndex == pageCount {
		return pageIndex
	}

	return pageIndex +1
}