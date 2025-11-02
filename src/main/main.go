package main

import (
	"fastlink/auth"
	"fastlink/config"
	"fastlink/internal"
	"fastlink/internal/link"
	"fastlink/internal/user"
	"fastlink/utils"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

// 查询计数中间件，每次请求自增 totalQuery
func queryCounterMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		atomic.AddUint64(&utils.TotalQuery, 1)
		ctx.Next()
	}
}

func main() {
	router := gin.Default()
	router.Use(queryCounterMiddleware())

	// 链接相关路由
	router.GET("/links/:url", func(ctx *gin.Context) {
		url := ctx.Param("url")
		link.GetLink(ctx, url)
	})
	router.DELETE("/links/:url", func(ctx *gin.Context) {
		url := ctx.Param("url")
		link.RemoveLink(ctx, url)
	})
	router.POST("/links/new", func(ctx *gin.Context) {
		link.GetShortUrl(ctx)
	})

	// 用户相关路由
	router.GET("/users/:account", func(ctx *gin.Context) {
		account := ctx.Param("account")
		user.GetUser(ctx, account)
	})
	router.DELETE("/users/:account", func(ctx *gin.Context) {
		account := ctx.Param("account")
		user.RemoveUser(ctx, account)
	})

	router.POST("/users/register", func(ctx *gin.Context) {
		auth.Register(ctx)
	})

	router.POST("/users/login", func(ctx *gin.Context) {
		auth.Login(ctx)
	})

	// 重定向路由
	router.GET("/:url", func(ctx *gin.Context) {
		url := ctx.Param("url")
		internal.Redirect(ctx, url)
	})

	// 启动 QPS 统计协程
	go utils.StatQps()

	router.Run(config.Global.Server.Port)
}
