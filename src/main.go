package main

import (
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

// 全局查询计数器
var totalQuery uint64 = 0

// 查询计数中间件，每次请求自增 totalQuery
func queryCounterMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		atomic.AddUint64(&totalQuery, 1)
		ctx.Next()
	}
}

func main() {
	router := gin.Default()
	router.Use(queryCounterMiddleware())

	// 链接相关路由
	router.GET("/links/:url", func(ctx *gin.Context) {
		url := ctx.Param("url")
		getLink(ctx, url)
	})
	router.DELETE("/links/:url", func(ctx *gin.Context) {
		url := ctx.Param("url")
		removeLink(ctx, url)
	})
	router.POST("/links/new", func(ctx *gin.Context) {
		getShortUrl(ctx)
	})

	// 用户相关路由
	router.GET("/users/:account", func(ctx *gin.Context) {
		account := ctx.Param("account")
		getUser(ctx, account)
	})
	router.DELETE("/users/:account", func(ctx *gin.Context) {
		account := ctx.Param("account")
		removeUser(ctx, account)
	})

	router.POST("/users/register", func(ctx *gin.Context) {
		register(ctx)
	})

	router.POST("/users/login", func(ctx *gin.Context) {
		login(ctx)
	})

	// 重定向路由
	router.GET("/:url", func(ctx *gin.Context) {
		url := ctx.Param("url")
		redirect(ctx, url)
	})

	// 启动 QPS 统计协程
	go statQps()

	router.Run()
}
