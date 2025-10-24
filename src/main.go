package main

import (
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

var totalQuery uint64 = 0

func queryCounterMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		atomic.AddUint64(&totalQuery, 1)
		ctx.Next()
	}
}

func main() {
	router := gin.Default()
	router.Use(queryCounterMiddleware())

	router.GET("/link/:url", func(ctx *gin.Context) {
		url := ctx.Param("url")
		getLink(ctx, url)
	})
	router.DELETE("/link/:url", func(ctx *gin.Context) {
		url := ctx.Param("url")
		removeLink(ctx, url)
	})
	router.POST("/link/new", func(ctx *gin.Context) {
		getShortUrl(ctx)
	})


	router.GET("/user/:account", func(ctx *gin.Context) {
		account := ctx.Param("account")
		getUser(ctx, account)
	})
	router.DELETE("/user/:account", func(ctx *gin.Context) {
		account := ctx.Param("account")
		removeUser(ctx, account)
	})

	router.POST("/user/register", func(ctx *gin.Context) {
		register(ctx)
	})

	router.POST("/user/login", func(ctx *gin.Context) {
		login(ctx)
	})

	

	router.GET("/:url", func(ctx *gin.Context) {
		url := ctx.Param("url")
		redirect(ctx, url)
	})

	//router group for rm user and link

	go statQps()

	router.Run()
}
