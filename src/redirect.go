package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

// redirect 根据短链重定向到原始链接，并增加点击数
func redirect(ctx *gin.Context, short string) {

	var link Link
	// 查询短链
	has, err := engine.Where("short_url=?", short).Get(&link)
	if err != nil {
		ctx.JSON(500, Response{Success: false, Error: "Database error"})
		return
	}
	if link.ExpireAt.Before(time.Now()) {
		ctx.JSON(404, Response{Success: false, Error: "Not found"})
		return
	}

	// 增加点击数
	_, err = engine.ID(link.ID).Update(&Link{ClickCount: link.ClickCount + 1})

	if err != nil {
		ctx.JSON(500, Response{Success: false, Error: "Internal server error"})
		return
	}

	if !has {
		ctx.JSON(404, Response{Success: false, Error: "Not found"})
		return
	}

	// 重定向到原始链接
	ctx.Redirect(301, link.SourceUrl)

}
