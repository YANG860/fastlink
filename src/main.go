package main

import (
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

var totalQuery uint64 = 0

func main() {
	router := gin.Default()

	router.POST("/register", func(ctx *gin.Context) {
		register(ctx)
		atomic.AddUint64(&totalQuery, 1)
	})

	router.POST("/login", func(ctx *gin.Context) {
		login(ctx)
		atomic.AddUint64(&totalQuery, 1)
	})

	router.POST("/new", func(ctx *gin.Context) {
		getShortUrl(ctx)
		atomic.AddUint64(&totalQuery, 1)
	})

	router.GET("/:url", func(ctx *gin.Context) {
		url := ctx.Param("url")
		redirect(ctx, url)

		atomic.AddUint64(&totalQuery, 1)
	})

	//router group for rm user and link

	removeGroup := router.Group("/remove")

	removeGroup.POST("/user", func(ctx *gin.Context) {
		removeUser(ctx)
		atomic.AddUint64(&totalQuery, 1)
	})

	removeGroup.POST("/link", func(ctx *gin.Context) {
		removeLink(ctx)
		atomic.AddUint64(&totalQuery, 1)
	})

	go statQps()

	router.Run()
}
