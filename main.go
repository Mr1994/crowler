package main

import (
	"api_client/config"
	"api_client/help"
	middleWare "api_client/middleware"
	"api_client/route"
	"github.com/gin-gonic/gin"
)

func main() {
	help.SetAppIniFile("main")
	// 初始化数据库
	config.Init()

	r := gin.Default()
	// 设置日志地址
	r.Use(middleWare.LogerMiddleware())
	// 设置请求
	r.Use(middleWare.RequestMiddleWare())
	// 设置路由器
	route.Routes(r)
	// 监听端口
	r.Run(":8001")
}

/**
切记 切记，如果某个结构体实现了这个某个inferface 那么就  这个结构体 = 这个infertface  b实现了c的方法，所以b可以当做inferface传入方法使用
*/
//type C interface {
//	Read()
//}
//
//// io render
//type A struct {
//	Name string
//}
//
//type B struct {
//	a A
//}
//
//func (b *B) Read()  {
//
//}
//
//func NewB(a A) *B {
//	return &B{
//		a: a,
//	}
//}
//
//func GetFormatB(c C)  {
//
//}
//
//func main() {
//	a := A{
//		Name: "zz",
//	}
//	b := NewB(a)
//
//	GetFormatB(b)
//}
